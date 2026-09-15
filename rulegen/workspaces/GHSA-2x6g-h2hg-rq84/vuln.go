package main

// Package manager contains plugin manager logic.
package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/grafana/grafana/pkg/infra/fs"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
	"github.com/grafana/grafana/pkg/plugins/manager/installer"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
	"github.com/grafana/grafana/pkg/util/errutil"
)

var (
	plog         log.Logger
	installerLog = NewInstallerLogger("plugin.installer", true)
)

const (
	grafanaComURL = "https://grafana.com/api/plugins"
)

type unsignedPluginConditionFunc = func(plugin *plugins.PluginBase) bool

type PluginScanner struct {
	pluginPath                    string
	errors                        []error
	backendPluginManager          backendplugin.Manager
	cfg                           *setting.Cfg
	requireSigned                 bool
	log                           log.Logger
	plugins                       map[string]*plugins.PluginBase
	allowUnsignedPluginsCondition unsignedPluginConditionFunc
}

type PluginManager struct {
	BackendPluginManager backendplugin.Manager `inject:""`
	Cfg                  *setting.Cfg          `inject:""`
	SQLStore             *sqlstore.SQLStore    `inject:""`
	pluginInstaller      plugins.PluginInstaller
	log                  log.Logger
	scanningErrors       []error

	// AllowUnsignedPluginsCondition changes the policy for allowing unsigned plugins. Signature validation only runs when plugins are starting
	// and running plugins will not be terminated if they violate the new policy.
	AllowUnsignedPluginsCondition unsignedPluginConditionFunc
	grafanaLatestVersion          string
	grafanaHasUpdate              bool
	pluginScanningErrors          map[string]plugins.PluginError

	renderer     *plugins.RendererPlugin
	dataSources  map[string]*plugins.DataSourcePlugin
	plugins      map[string]*plugins.PluginBase
	panels       map[string]*plugins.PanelPlugin
	apps         map[string]*plugins.AppPlugin
	staticRoutes []*plugins.PluginStaticRoute
	pluginsMu    sync.RWMutex
}

func init() {
	registry.Register(&registry.Descriptor{
		Name:         "PluginManager",
		Instance:     newManager(nil),
		InitPriority: registry.MediumHigh,
	})
}

func newManager(cfg *setting.Cfg) *PluginManager {
	return &PluginManager{
		Cfg:         cfg,
		dataSources: map[string]*plugins.DataSourcePlugin{},
		plugins:     map[string]*plugins.PluginBase{},
		panels:      map[string]*plugins.PanelPlugin{},
		apps:        map[string]*plugins.AppPlugin{},
	}
}

func (pm *PluginManager) Init() error {
	pm.log = log.New("plugins")
	plog = log.New("plugins")
	pm.pluginScanningErrors = map[string]plugins.PluginError{}
	pm.pluginInstaller = installer.New(false, pm.Cfg.BuildVersion, installerLog)

	pm.log.Info("Starting plugin search")

	plugDir := filepath.Join(pm.Cfg.StaticRootPath, "app/plugins")
	pm.log.Debug("Scanning core plugin directory", "dir", plugDir)
	if err := pm.scan(plugDir, false); err != nil {
		return errutil.Wrapf(err, "failed to scan core plugin directory '%s'", plugDir)
	}

	plugDir = pm.Cfg.BundledPluginsPath
	pm.log.Debug("Scanning bundled plugins directory", "dir", plugDir)
	exists, err := fs.Exists(plugDir)
	if err != nil {
		return err
	}
	if exists {
		if err := pm.scan(plugDir, false); err != nil {
			return errutil.Wrapf(err, "failed to scan bundled plugins directory '%s'", plugDir)
		}
	}

	err = pm.initExternalPlugins()
	if err != nil {
		return err
	}

	return nil
}

func (pm *PluginManager) initExternalPlugins() error {
	// check if plugins dir exists
	exists, err := fs.Exists(pm.Cfg.PluginsPath)
	if err != nil {
		return err
	}

	if !exists {
		if err = os.MkdirAll(pm.Cfg.PluginsPath, os.ModePerm); err != nil {
			pm.log.Error("failed to create external plugins directory", "dir", pm.Cfg.PluginsPath, "error", err)
		} else {
			pm.log.Info("External plugins directory created", "directory", pm.Cfg.PluginsPath)
		}
	} else {
		pm.log.Debug("Scanning external plugins directory", "dir", pm.Cfg.PluginsPath)
		if err := pm.scan(pm.Cfg.PluginsPath, true); err != nil {
			return errutil.Wrapf(err, "failed to scan external plugins directory '%s'",
				pm.Cfg.PluginsPath)
		}
	}

	if err := pm.scanPluginPaths(); err != nil {
		return err
	}

	var staticRoutesList []*plugins.PluginStaticRoute
	for _, panel := range pm.Panels() {
		staticRoutes := panel.InitFrontendPlugin(pm.Cfg)
		staticRoutesList = append(staticRoutesList, staticRoutes...)
	}

	for _, ds := range pm.DataSources() {
		staticRoutes := ds.InitFrontendPlugin(pm.Cfg)
		staticRoutesList = append(staticRoutesList, staticRoutes...)
	}

	for _, app := range pm.Apps() {
		staticRoutes := app.InitApp(pm.panels, pm.dataSources, pm.Cfg)
		staticRoutesList = append(staticRoutesList, staticRoutes...)
	}

	if pm.Renderer() != nil {
		staticRoutes := pm.renderer.InitFrontendPlugin(pm.Cfg)
		staticRoutesList = append(staticRoutesList, staticRoutes...)
	}
	pm.staticRoutes = staticRoutesList

	for _, p := range pm.Plugins() {
		if p.IsCorePlugin {
			p.Signature = plugins.PluginSignatureInternal
		} else {
			metrics.SetPluginBuildInformation(p.Id, p.Type, p.Info.Version)
		}
	}

	return nil
}

func (pm *PluginManager) Run(ctx context.Context) error {
	pm.checkForUpdates()

	ticker := time.NewTicker(time.Minute * 10)
	run := true

	for run {
		select {
		case <-ticker.C:
			pm.checkForUpdates()
		case <-ctx.Done():
			run = false
		}
	}

	return ctx.Err()
}

func (pm *PluginManager) Renderer() *plugins.RendererPlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return pm.renderer
}

func (pm *PluginManager) GetDataSource(id string) *plugins.DataSourcePlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return pm.dataSources[id]
}

func (pm *PluginManager) DataSources() []*plugins.DataSourcePlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	var rslt []*plugins.DataSourcePlugin
	for _, ds := range pm.dataSources {
		rslt = append(rslt, ds)
	}

	return rslt
}

func (pm *PluginManager) DataSourceCount() int {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return len(pm.dataSources)
}

func (pm *PluginManager) PanelCount() int {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return len(pm.panels)
}

func (pm *PluginManager) AppCount() int {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return len(pm.apps)
}

func (pm *PluginManager) Plugins() []*plugins.PluginBase {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	var rslt []*plugins.PluginBase
	for _, p := range pm.plugins {
		rslt = append(rslt, p)
	}

	return rslt
}

func (pm *PluginManager) Apps() []*plugins.AppPlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	var rslt []*plugins.AppPlugin
	for _, p := range pm.apps {
		rslt = append(rslt, p)
	}

	return rslt
}

func (pm *PluginManager) Panels() []*plugins.PanelPlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	var rslt []*plugins.PanelPlugin
	for _, p := range pm.panels {
		rslt = append(rslt, p)
	}

	return rslt
}

func (pm *PluginManager) GetPlugin(id string) *plugins.PluginBase {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return pm.plugins[id]
}

func (pm *PluginManager) GetApp(id string) *plugins.AppPlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	return pm.apps[id]
}

func (pm *PluginManager) GrafanaLatestVersion() string {
	return pm.grafanaLatestVersion
}

func (pm *PluginManager) GrafanaHasUpdate() bool {
	return pm.grafanaHasUpdate
}

// scanPluginPaths scans configured plugin paths.
func (pm *PluginManager) scanPluginPaths() error {
	for pluginID, settings := range pm.Cfg.PluginSettings {
		path, exists := settings["path"]
		if !exists || path == "" {
			continue
		}

		if err := pm.scan(path, true); err != nil {
			return errutil.Wrapf(err, "failed to scan directory configured for plugin '%s': '%s'", pluginID, path)
		}
	}

	return nil
}

// scan a directory for plugins.
func (pm *PluginManager) scan(pluginDir string, requireSigned bool) error {
	scanner := &PluginScanner{
		pluginPath:                    pluginDir,
		backendPluginManager:          pm.BackendPluginManager,
		cfg:                           pm.Cfg,
		requireSigned:                 requireSigned,
		log:                           pm.log,
		plugins:                       map[string]*plugins.PluginBase{},
		allowUnsignedPluginsCondition: pm.AllowUnsignedPluginsCondition,
	}

	// 1st pass: Scan plugins, also mapping plugins to their respective directories
	if err := util.Walk(pluginDir, true, true, scanner.walker); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			pm.log.Debug("Couldn't scan directory since it doesn't exist", "pluginDir", pluginDir, "err", err)
			return nil
		}
		if errors.Is(err, os.ErrPermission) {
			pm.log.Debug("Couldn't scan directory due to lack of permissions", "pluginDir", pluginDir, "err", err)
			return nil
		}
		if pluginDir != "data/plugins" {
			pm.log.Warn("Could not scan dir", "pluginDir", pluginDir, "err", err)
		}
		return err
	}

	pm.log.Debug("Initial plugin loading done")

	pluginsByID := make(map[string]struct{})
	for scannedPluginPath, scannedPlugin := range scanner.plugins {
		// Check if scanning found duplicate plugins
		if _, dupe := pluginsByID[scannedPlugin.Id]; dupe {
			pm.log.Warn("Skipping plugin as it's a duplicate", "id", scannedPlugin.Id)
			scanner.errors = append(scanner.errors,
				plugins.DuplicatePluginError{PluginID: scannedPlugin.Id, ExistingPluginDir: scannedPlugin.PluginDir})
			delete(scanner.plugins, scannedPluginPath)
			continue
		}
		pluginsByID[scannedPlugin.Id] = struct{}{}

		// Check if scanning found plugins that are already installed
		if existing := pm.GetPlugin(scannedPlugin.Id); existing != nil {
			pm.log.Debug("Skipping plugin as it's already installed", "plugin", existing.Id, "version", existing.Info.Version)
			delete(scanner.plugins, scannedPluginPath)
		}
	}

	pluginTypes := map[string]interface{}{
		"panel":      plugins.PanelPlugin{},
		"datasource": plugins.DataSourcePlugin{},
		"app":        plugins.AppPlugin{},
		"renderer":   plugins.RendererPlugin{},
	}

	// 2nd pass: Validate and register plugins
	for dpath, plugin := range scanner.plugins {
		// Try to find any root plugin
		ancestors := strings.Split(dpath, string(filepath.Separator))
		ancestors = ancestors[0 : len(ancestors)-1]
		aPath := ""
		if runtime.GOOS != "windows" && filepath.IsAbs(dpath) {
			aPath = "/"
		}
		for _, a := range ancestors {
			aPath = filepath.Join(aPath, a)
			if root, ok := scanner.plugins[aPath]; ok {
				plugin.Root = root
				break
			}
		}

		pm.log.Debug("Found plugin", "id", plugin.Id, "signature", plugin.Signature, "hasRoot", plugin.Root != nil)
		signingError := scanner.validateSignature(plugin)
		if signingError != nil {
			pm.log.Debug("Failed to validate plugin signature. Will skip loading", "id", plugin.Id,
				"signature", plugin.Signature, "status", signingError.ErrorCode)
			pm.pluginScanningErrors[plugin.Id] = *signingError
			continue
		}

		pm.log.Debug("Attempting to add plugin", "id", plugin.Id)

		pluginGoType, exists := pluginTypes[plugin.Type]
		if !exists {
			return fmt.Errorf("unknown plugin type %q", plugin.Type)
		}

		jsonFPath := filepath.Join(plugin.PluginDir, "plugin.json")

		// External plugins need a module.js file for SystemJS to load
		if !strings.HasPrefix(jsonFPath, pm.Cfg.StaticRootPath) && !scanner.IsBackendOnlyPlugin(plugin.Type) {
			module := filepath.Join(plugin.PluginDir, "module.js")
			exists, err := fs.Exists(module)
			if err != nil {
				return err
			}
			if !exists {
				scanner.log.Warn("Plugin missing module.js",
					"name", plugin.Name,
					"warning", "Missing module.js, If you loaded this plugin from git, make sure to compile it.",
					"path", module)
			}
		}

		// nolint:gosec
		// We can ignore the gosec G304 warning on this one because `jsonFPath` is based
		// on plugin the folder structure on disk and not user input.
		reader, err := os.Open(jsonFPath)
		if err != nil {
			return err
		}
		defer func() {
			if err := reader.Close(); err != nil {
				scanner.log.Warn("Failed to close JSON file", "path", jsonFPath, "err", err)
			}
		}()

		jsonParser := json.NewDecoder(reader)

		loader := reflect.New(reflect.TypeOf(pluginGoType)).Interface().(plugins.PluginLoader)

		// Load the full plugin, and add it to manager
		if err := pm.loadPlugin(jsonParser, plugin, scanner, loader); err != nil {
			return err
		}
	}

	if len(scanner.errors) > 0 {
		var errStr []string
		for _, err := range scanner.errors {
			errStr = append(errStr, err.Error())
		}
		pm.log.Warn("Some plugin scanning errors were found", "errors", strings.Join(errStr, ", "))
		pm.scanningErrors = scanner.errors
	}

	return nil
}

func (pm *PluginManager) loadPlugin(jsonParser *json.Decoder, pluginBase *plugins.PluginBase,
	scanner *PluginScanner, loader plugins.PluginLoader) error {
	plug, err := loader.Load(jsonParser, pluginBase, scanner.backendPluginManager)
	if err != nil {
		return err
	}

	pm.pluginsMu.Lock()
	defer pm.pluginsMu.Unlock()

	var pb *plugins.PluginBase
	switch p := plug.(type) {
	case *plugins.DataSourcePlugin:
		pm.dataSources[p.Id] = p
		pb = &p.PluginBase
	case *plugins.PanelPlugin:
		pm.panels[p.Id] = p
		pb = &p.PluginBase
	case *plugins.RendererPlugin:
		pm.renderer = p
		pb = &p.PluginBase
	case *plugins.AppPlugin:
		pm.apps[p.Id] = p
		pb = &p.PluginBase
	default:
		panic(fmt.Sprintf("Unrecognized plugin type %T", plug))
	}

	if !strings.HasPrefix(pluginBase.PluginDir, pm.Cfg.StaticRootPath) {
		pm.log.Info("Registering plugin", "id", pb.Id)
	}

	if len(pb.Dependencies.Plugins) == 0 {
		pb.Dependencies.Plugins = []plugins.PluginDependencyItem{}
	}

	if pb.Dependencies.GrafanaVersion == "" {
		pb.Dependencies.GrafanaVersion = "*"
	}

	for _, include := range pb.Includes {
		if include.Role == "" {
			include.Role = models.ROLE_VIEWER
		}
	}

	// Copy relevant fields from the base
	pb.PluginDir = pluginBase.PluginDir
	pb.Signature = pluginBase.Signature
	pb.SignatureType = pluginBase.SignatureType
	pb.SignatureOrg = pluginBase.SignatureOrg

	pm.plugins[pb.Id] = pb
	pm.log.Debug("Successfully added plugin", "id", pb.Id)
	return nil
}

func (s *PluginScanner) walker(currentPath string, f os.FileInfo, err error) error {
	// We scan all the subfolders for plugin.json (with some exceptions) so that we also load embedded plugins, for
	// example https://github.com/raintank/worldping-app/tree/master/dist/grafana-worldmap-panel worldmap panel plugin
	// is embedded in worldping app.
	if err != nil {
		return fmt.Errorf("filepath.Walk reported an error for %q: %w", currentPath, err)
	}

	if f.Name() == "node_modules" || f.Name() == "Chromium.app" {
		return util.ErrWalkSkipDir
	}

	if f.IsDir() {
		return nil
	}

	if f.Name() != "plugin.json" {
		return nil
	}

	if err := s.loadPlugin(currentPath); err != nil {
		s.log.Error("Failed to load plugin", "error", err, "pluginPath", filepath.Dir(currentPath))
		s.errors = append(s.errors, err)
	}

	return nil
}

func (s *PluginScanner) loadPlugin(pluginJSONFilePath string) error {
	s.log.Debug("Loading plugin", "path", pluginJSONFilePath)
	currentDir := filepath.Dir(pluginJSONFilePath)
	// nolint:gosec
	// We can ignore the gosec G304 warning on this one because `currentPath` is based
	// on plugin the folder structure on disk and not user input.
	reader, err := os.Open(pluginJSONFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := reader.Close(); err != nil {
			s.log.Warn("Failed to close JSON file", "path", pluginJSONFilePath, "err", err)
		}
	}()

	jsonParser := json.NewDecoder(reader)
	pluginCommon := plugins.PluginBase{}
	if err := jsonParser.Decode(&pluginCommon); err != nil {
		return err
	}

	if pluginCommon.Id == "" || pluginCommon.Type == "" {
		return errors.New("did not find type or id properties in plugin.json")
	}

	pluginCommon.PluginDir = filepath.Dir(pluginJSONFilePath)
	pluginCommon.Files, err = collectPluginFilesWithin(pluginCommon.PluginDir)
	if err != nil {
		s.log.Warn("Could not collect plugin file information in directory", "pluginID", pluginCommon.Id, "dir", pluginCommon.PluginDir)
		return err
	}

	signatureState, err := getPluginSignatureState(s.log, &pluginCommon)
	if err != nil {
		s.log.Warn("Could not get plugin signature state", "pluginID", pluginCommon.Id, "err", err)
		return err
	}
	pluginCommon.Signature = signatureState.Status
	pluginCommon.SignatureType = signatureState.Type
	pluginCommon.SignatureOrg = signatureState.SigningOrg

	s.plugins[currentDir] = &pluginCommon

	return nil
}

func (*PluginScanner) IsBackendOnlyPlugin(pluginType string) bool {
	return pluginType == "renderer"
}

// validateSignature validates a plugin's signature.
func (s *PluginScanner) validateSignature(plugin *plugins.PluginBase) *plugins.PluginError {
	if plugin.Signature == plugins.PluginSignatureValid {
		s.log.Debug("Plugin has valid signature", "id", plugin.Id)
		return nil
	}

	if plugin.Root != nil {
		// If a descendant plugin with invalid signature, set signature to that of root
		if plugin.IsCorePlugin || plugin.Signature == plugins.PluginSignatureInternal {
			s.log.Debug("Not setting descendant plugin's signature to that of root since it's core or internal",
				"plugin", plugin.Id, "signature", plugin.Signature, "isCore", plugin.IsCorePlugin)
		} else {
			s.log.Debug("Setting descendant plugin's signature to that of root", "plugin", plugin.Id,
				"root", plugin.Root.Id, "signature", plugin.Signature, "rootSignature", plugin.Root.Signature)
			plugin.Signature = plugin.Root.Signature
			if plugin.Signature == plugins.PluginSignatureValid {
				s.log.Debug("Plugin has valid signature (inherited from root)", "id", plugin.Id)
				return nil
			}
		}
	} else {
		s.log.Debug("Non-valid plugin Signature", "pluginID", plugin.Id, "pluginDir", plugin.PluginDir,
			"state", plugin.Signature)
	}

	if !s.requireSigned {
		return nil
	}

	switch plugin.Signature {
	case plugins.PluginSignatureUnsigned:
		if allowed := s.allowUnsigned(plugin); !allowed {
			s.log.Debug("Plugin is unsigned", "id", plugin.Id)
			s.errors = append(s.errors, fmt.Errorf("plugin '%s' is unsigned", plugin.Id))
			return &plugins.PluginError{
				ErrorCode: signatureMissing,
			}
		}
		s.log.Warn("Running an unsigned plugin", "pluginID", plugin.Id, "pluginDir",
			plugin.PluginDir)
		return nil
	case plugins.PluginSignatureInvalid:
		s.log.Debug("Plugin '%s' has an invalid signature", plugin.Id)
		s.errors = append(s.errors, fmt.Errorf("plugin '%s' has an invalid signature", plugin.Id))
		return &plugins.PluginError{
			ErrorCode: signatureInvalid,
		}
	case plugins.PluginSignatureModified:
		s.log.Debug("Plugin '%s' has a modified signature", plugin.Id)
		s.errors = append(s.errors, fmt.Errorf("plugin '%s' has a modified signature", plugin.Id))
		return &plugins.PluginError{
			ErrorCode: signatureModified,
		}
	default:
		panic(fmt.Sprintf("Plugin '%s' has an unrecognized plugin signature state '%s'", plugin.Id, plugin.Signature))
	}
}

func (s *PluginScanner) allowUnsigned(plugin *plugins.PluginBase) bool {
	if s.allowUnsignedPluginsCondition != nil {
		return s.allowUnsignedPluginsCondition(plugin)
	}

	if s.cfg.Env == setting.Dev {
		return true
	}

	for _, plug := range s.cfg.PluginsAllowUnsigned {
		if plug == plugin.Id {
			return true
		}
	}

	return false
}

// ScanningErrors returns plugin scanning errors encountered.
func (pm *PluginManager) ScanningErrors() []plugins.PluginError {
	scanningErrs := make([]plugins.PluginError, 0)
	for id, e := range pm.pluginScanningErrors {
		scanningErrs = append(scanningErrs, plugins.PluginError{
			ErrorCode: e.ErrorCode,
			PluginID:  id,
		})
	}
	return scanningErrs
}

func (pm *PluginManager) GetPluginMarkdown(pluginId string, name string) ([]byte, error) {
	plug, exists := pm.plugins[pluginId]
	if !exists {
		return nil, plugins.PluginNotFoundError{PluginID: pluginId}
	}

	// nolint:gosec
	// We can ignore the gosec G304 warning on this one because `plug.PluginDir` is based
	// on plugin the folder structure on disk and not user input.
	path := filepath.Join(plug.PluginDir, fmt.Sprintf("%s.md", strings.ToUpper(name)))
	exists, err := fs.Exists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		path = filepath.Join(plug.PluginDir, fmt.Sprintf("%s.md", strings.ToLower(name)))
	}

	exists, err = fs.Exists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return make([]byte, 0), nil
	}

	// nolint:gosec
	// We can ignore the gosec G304 warning on this one because `plug.PluginDir` is based
	// on plugin the folder structure on disk and not user input.
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// gets plugin filenames that require verification for plugin signing
func collectPluginFilesWithin(rootDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() != "MANIFEST.txt" {
			file, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(file))
		}
		return nil
	})
	return files, err
}

// GetDataPlugin gets a DataPlugin with a certain name. If none is found, nil is returned.
//nolint: staticcheck // plugins.DataPlugin deprecated
func (pm *PluginManager) GetDataPlugin(id string) plugins.DataPlugin {
	pm.pluginsMu.RLock()
	defer pm.pluginsMu.RUnlock()

	if p := pm.GetDataSource(id); p != nil && p.CanHandleDataQueries() {
		return p
	}

	// XXX: Might other plugins implement DataPlugin?

	p := pm.BackendPluginManager.GetDataPlugin(id)
	if p != nil {
		return p.(plugins.DataPlugin)
	}

	return nil
}

func (pm *PluginManager) StaticRoutes() []*plugins.PluginStaticRoute {
	return pm.staticRoutes
}

func (pm *PluginManager) Install(ctx context.Context, pluginID, version string) error {
	plugin := pm.GetPlugin(pluginID)
	if plugin != nil {
		if plugin.IsCorePlugin {
			return plugins.ErrInstallCorePlugin
		}

		if plugin.Info.Version == version {
			return plugins.DuplicatePluginError{
				PluginID:          pluginID,
				ExistingPluginDir: plugin.PluginDir,
			}
		}

		// remove existing installation of plugin
		err := pm.Uninstall(context.Background(), plugin.Id)
		if err != nil {
			return err
		}
	}

	err := pm.pluginInstaller.Install(ctx, pluginID, version, pm.Cfg.PluginsPath, "", grafanaComURL)
	if err != nil {
		return err
	}

	err = pm.initExternalPlugins()
	if err != nil {
		return err
	}

	return nil
}

func (pm *PluginManager) Uninstall(ctx context.Context, pluginID string) error {
	plugin := pm.GetPlugin(pluginID)
	if plugin == nil {
		return plugins.ErrPluginNotInstalled
	}

	if plugin.IsCorePlugin {
		return plugins.ErrUninstallCorePlugin
	}

	// extra security check to ensure we only remove plugins that are located in the configured plugins directory
	path, err := filepath.Rel(pm.Cfg.PluginsPath, plugin.PluginDir)
	if err != nil || strings.HasPrefix(path, ".."+string(filepath.Separator)) {
		return plugins.ErrUninstallOutsideOfPluginDir
	}

	if pm.BackendPluginManager.IsRegistered(pluginID) {
		err := pm.BackendPluginManager.UnregisterAndStop(ctx, pluginID)
		if err != nil {
			return err
		}
	}

	err = pm.unregister(plugin)
	if err != nil {
		return err
	}

	return pm.pluginInstaller.Uninstall(ctx, pluginID, pm.Cfg.PluginsPath)
}

func (pm *PluginManager) unregister(plugin *plugins.PluginBase) error {
	pm.pluginsMu.Lock()
	defer pm.pluginsMu.Unlock()

	switch plugin.Type {
	case "panel":
		delete(pm.panels, plugin.Id)
	case "datasource":
		delete(pm.dataSources, plugin.Id)
	case "app":
		delete(pm.apps, plugin.Id)
	case "renderer":
		pm.renderer = nil
	}

	delete(pm.plugins, plugin.Id)

	pm.removeStaticRoute(plugin.Id)

	return nil
}

func (pm *PluginManager) removeStaticRoute(pluginID string) {
	for i, route := range pm.staticRoutes {
		if pluginID == route.PluginId {
			pm.staticRoutes = append(pm.staticRoutes[:i], pm.staticRoutes[i+1:]...)
			return
		}
	}
}
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
	"github.com/grafana/grafana/pkg/plugins/manager/installer"
	"github.com/grafana/grafana/pkg/setting"
)

func (hs *HTTPServer) GetPluginList(c *models.ReqContext) response.Response {
		coreFilter = "1"
	}

	pluginSettingsMap, err := hs.PluginManager.GetPluginSettings(c.OrgId)
	if err != nil {
		return response.Error(500, "Failed to get list of plugins", err)
	}

	result := make(dtos.PluginList, 0)
	for _, pluginDef := range hs.PluginManager.Plugins() {
		// filter out app sub plugins
		if embeddedFilter == "0" && pluginDef.IncludedInAppId != "" {
			continue
		}

		// filter out core plugins
		if (coreFilter == "0" && pluginDef.IsCorePlugin) || (coreFilter == "1" && !pluginDef.IsCorePlugin) {
			continue
		}

		// filter on type
		if typeFilter != "" && typeFilter != pluginDef.Type {
			continue
		}

		if pluginDef.State == plugins.PluginStateAlpha && !hs.Cfg.PluginsEnableAlpha {
			continue
		}

		listItem := dtos.PluginListItem{
			Id:            pluginDef.Id,
			Name:          pluginDef.Name,
			Type:          pluginDef.Type,
			Category:      pluginDef.Category,
			Info:          &pluginDef.Info,
			LatestVersion: pluginDef.GrafanaNetVersion,
			HasUpdate:     pluginDef.GrafanaNetHasUpdate,
			DefaultNavUrl: pluginDef.DefaultNavUrl,
			State:         pluginDef.State,
			Signature:     pluginDef.Signature,
			SignatureType: pluginDef.SignatureType,
			SignatureOrg:  pluginDef.SignatureOrg,
		}

		if pluginSetting, exists := pluginSettingsMap[pluginDef.Id]; exists {
			listItem.Enabled = pluginSetting.Enabled
			listItem.Pinned = pluginSetting.Pinned
		}
			continue
		}

		// filter out built in data sources
		if ds := hs.PluginManager.GetDataSource(pluginDef.Id); ds != nil {
			if ds.BuiltIn {
				continue
			}
		}

		result = append(result, listItem)
}

func (hs *HTTPServer) GetPluginSettingByID(c *models.ReqContext) response.Response {
	pluginID := c.Params(":pluginId")

	def := hs.PluginManager.GetPlugin(pluginID)
	if def == nil {
		return response.Error(404, "Plugin not found, no installed plugin with that id", nil)
	}

	dto := &dtos.PluginSetting{
		Type:          def.Type,
		Id:            def.Id,
		Name:          def.Name,
		Info:          &def.Info,
		Dependencies:  &def.Dependencies,
		Includes:      def.Includes,
		BaseUrl:       def.BaseUrl,
		Module:        def.Module,
		DefaultNavUrl: def.DefaultNavUrl,
		LatestVersion: def.GrafanaNetVersion,
		HasUpdate:     def.GrafanaNetHasUpdate,
		State:         def.State,
		Signature:     def.Signature,
		SignatureType: def.SignatureType,
		SignatureOrg:  def.SignatureOrg,
	}

	if app := hs.PluginManager.GetApp(def.Id); app != nil {
		dto.Enabled = app.AutoEnabled
		dto.Pinned = app.AutoEnabled
	}

	query := models.GetPluginSettingByIdQuery{PluginId: pluginID, OrgId: c.OrgId}
	if err := bus.Dispatch(&query); err != nil {
		if !errors.Is(err, models.ErrPluginSettingNotFound) {
			return response.Error(500, "Failed to get login settings", nil)
		}
	} else {
		dto.Enabled = query.Result.Enabled
		dto.Pinned = query.Result.Pinned
		dto.JsonData = query.Result.JsonData
	}

	return response.JSON(200, dto)
}

func (hs *HTTPServer) UpdatePluginSetting(c *models.ReqContext, cmd models.UpdatePluginSettingCmd) response.Response {
	pluginID := c.Params(":pluginId")

	if app := hs.PluginManager.GetApp(pluginID); app == nil {
		return response.Error(404, "Plugin not installed", nil)
	}

	cmd.OrgId = c.OrgId
	cmd.PluginId = pluginID
	if err := bus.Dispatch(&cmd); err != nil {
		return response.Error(500, "Failed to update plugin setting", err)
	}

	return response.Success("Plugin settings updated")
}

func (hs *HTTPServer) GetPluginDashboards(c *models.ReqContext) response.Response {
	pluginID := c.Params(":pluginId")

	list, err := hs.PluginManager.GetPluginDashboards(c.OrgId, pluginID)
	if err != nil {
		var notFound plugins.PluginNotFoundError
		if errors.As(err, &notFound) {
			return response.Error(404, notFound.Error(), nil)
		}

		return response.Error(500, "Failed to get plugin dashboards", err)
	}

	return response.JSON(200, list)
}

func (hs *HTTPServer) GetPluginMarkdown(c *models.ReqContext) response.Response {
	pluginID := c.Params(":pluginId")
	name := c.Params(":name")

	content, err := hs.PluginManager.GetPluginMarkdown(pluginID, name)
	if err != nil {
		var notFound plugins.PluginNotFoundError
		if errors.As(err, &notFound) {
			return response.Error(404, notFound.Error(), nil)
		}

	// fallback try readme
	if len(content) == 0 {
		content, err = hs.PluginManager.GetPluginMarkdown(pluginID, "readme")
		if err != nil {
			return response.Error(501, "Could not get markdown file", err)
		}
	return resp
}

func (hs *HTTPServer) ImportDashboard(c *models.ReqContext, apiCmd dtos.ImportDashboardCommand) response.Response {
	var err error
	if apiCmd.PluginId == "" && apiCmd.Dashboard == nil {
		return response.Error(422, "Dashboard must be set", nil)
	}

	trimDefaults := c.QueryBoolWithDefault("trimdefaults", true)
	if trimDefaults && !hs.LoadSchemaService.IsDisabled() {
		apiCmd.Dashboard, err = hs.LoadSchemaService.DashboardApplyDefaults(apiCmd.Dashboard)
		if err != nil {
			return response.Error(500, "Error while applying default value to the dashboard json", err)
		}
	}

	dashInfo, err := hs.PluginManager.ImportDashboard(apiCmd.PluginId, apiCmd.Path, c.OrgId, apiCmd.FolderId,
		apiCmd.Dashboard, apiCmd.Overwrite, apiCmd.Inputs, c.SignedInUser, hs.DataService)
	if err != nil {
		return hs.dashboardSaveErrorToApiResponse(err)
	}

	err = hs.LibraryPanelService.ImportDashboard(c, apiCmd.Dashboard, dashInfo.DashboardId)
	if err != nil {
		return response.Error(500, "Error while connecting library panels", err)
	}

	return response.JSON(200, dashInfo)
}

// CollectPluginMetrics collect metrics from a plugin.
//
// /api/plugins/:pluginId/metrics
func (hs *HTTPServer) CollectPluginMetrics(c *models.ReqContext) response.Response {
	pluginID := c.Params("pluginId")
	plugin := hs.PluginManager.GetPlugin(pluginID)
	if plugin == nil {
		return response.Error(404, "Plugin not found", nil)
	}

	resp, err := hs.BackendPluginManager.CollectMetrics(c.Req.Context(), plugin.Id)
	if err != nil {
		return translatePluginRequestErrorToAPIError(err)
	}
	return response.CreateNormalResponse(headers, resp.PrometheusMetrics, http.StatusOK)
}

// GetPluginAssets returns public plugin assets (images, JS, etc.)
//
// /public/plugins/:pluginId/*
func (hs *HTTPServer) GetPluginAssets(c *models.ReqContext) {
	pluginID := c.Params("pluginId")
	plugin := hs.PluginManager.GetPlugin(pluginID)
	if plugin == nil {
		c.Handle(hs.Cfg, 404, "Plugin not found", nil)
		return
	}

	requestedFile := filepath.Clean(c.Params("*"))
	pluginFilePath := filepath.Join(plugin.PluginDir, requestedFile)

	// It's safe to ignore gosec warning G304 since we already clean the requested file path and subsequently
	// use this with a prefix of the plugin's directory, which is set during plugin loading
	// nolint:gosec
	f, err := os.Open(pluginFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.Handle(hs.Cfg, 404, "Could not find plugin file", err)
			return
		}
		c.Handle(hs.Cfg, 500, "Could not open plugin file", err)
		return
	}
	defer func() {

	fi, err := f.Stat()
	if err != nil {
		c.Handle(hs.Cfg, 500, "Plugin file exists but could not open", err)
		return
	}

	if shouldExclude(fi) {
		c.Handle(hs.Cfg, 404, "Plugin file not found", nil)
		return
	}

	headers := func(c *macaron.Context) {
		c.Resp.Header().Set("Cache-Control", "public, max-age=3600")
	}

	if hs.Cfg.Env == setting.Dev {
		headers = func(c *macaron.Context) {
			c.Resp.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache")
		}
	}

	headers(c.Context)

	http.ServeContent(c.Resp, c.Req.Request, pluginFilePath, fi.ModTime(), f)
}

// CheckHealth returns the health of a plugin.
// /api/plugins/:pluginId/health
func (hs *HTTPServer) CheckHealth(c *models.ReqContext) response.Response {
	pluginID := c.Params("pluginId")

	pCtx, found, err := hs.PluginContextProvider.Get(pluginID, "", c.SignedInUser, false)
	if err != nil {
		return response.Error(500, "Failed to get plugin settings", err)
	}
		return response.Error(404, "Plugin not found", nil)
	}

	resp, err := hs.BackendPluginManager.CheckHealth(c.Req.Context(), pCtx)
	if err != nil {
		return translatePluginRequestErrorToAPIError(err)
	}
//
// /api/plugins/:pluginId/resources/*
func (hs *HTTPServer) CallResource(c *models.ReqContext) {
	pluginID := c.Params("pluginId")

	pCtx, found, err := hs.PluginContextProvider.Get(pluginID, "", c.SignedInUser, false)
	if err != nil {
		c.JsonApiErr(500, "Failed to get plugin settings", err)
		return
	}
	if !found {
		c.JsonApiErr(404, "Plugin not found", nil)
		return
	}
	hs.BackendPluginManager.CallResource(pCtx, c, c.Params("*"))
}

func (hs *HTTPServer) GetPluginErrorsList(_ *models.ReqContext) response.Response {
	return response.JSON(200, hs.PluginManager.ScanningErrors())
}

func (hs *HTTPServer) InstallPlugin(c *models.ReqContext, dto dtos.InstallPluginCommand) response.Response {
	pluginID := c.Params("pluginId")

	err := hs.PluginManager.Install(c.Req.Context(), pluginID, dto.Version)
	if err != nil {
		var dupeErr plugins.DuplicatePluginError
		if errors.As(err, &dupeErr) {
			return response.Error(http.StatusConflict, "Plugin already installed", err)
		}
		if errors.As(err, &versionNotFoundErr) {
			return response.Error(http.StatusNotFound, "Plugin version not found", err)
		}
		if errors.Is(err, installer.ErrPluginNotFound) {
			return response.Error(http.StatusNotFound, "Plugin not found", err)
		}
		if errors.Is(err, plugins.ErrInstallCorePlugin) {
			return response.Error(http.StatusForbidden, "Cannot install or change a Core plugin", err)
}

func (hs *HTTPServer) UninstallPlugin(c *models.ReqContext) response.Response {
	pluginID := c.Params("pluginId")

	err := hs.PluginManager.Uninstall(c.Req.Context(), pluginID)
	if err != nil {
		if errors.Is(err, plugins.ErrPluginNotInstalled) {
			return response.Error(http.StatusNotFound, "Plugin not installed", err)
	return response.Error(500, "Plugin request failed", err)
}

func shouldExclude(fi os.FileInfo) bool {
	normalizedFilename := strings.ToLower(fi.Name())

	isUnixExecutable := fi.Mode()&0111 == 0111
	isWindowsExecutable := strings.HasSuffix(normalizedFilename, ".exe")
	isScript := strings.HasSuffix(normalizedFilename, ".sh")

	return isUnixExecutable || isWindowsExecutable || isScript
}

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/setting"
)

		},
	}

	if setting.AddChangePasswordLink() {
		children = append(children, &dtos.NavLink{
			Text: "Change password", Id: "change-password", Url: hs.Cfg.AppSubURL + "/profile/password",
			Icon: "lock", HideFromMenu: true,
		})
	}

	}

	return &dtos.NavLink{
		Text:         c.SignedInUser.NameOrFallback(),
		SubTitle:     login,
		Id:           "profile",
		Img:          gravatarURL,
		Url:          hs.Cfg.AppSubURL + "/profile",
		HideFromMenu: true,
		SortWeight:   dtos.WeightProfile,
		Children:     children,
	}
}

func (hs *HTTPServer) getAppLinks(c *models.ReqContext) ([]*dtos.NavLink, error) {
	enabledPlugins, err := hs.PluginManager.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}

	appLinks := []*dtos.NavLink{}
	for _, plugin := range enabledPlugins.Apps {
		if !plugin.Pinned {
			continue
		}

		appLink := &dtos.NavLink{
			Text:       plugin.Name,
			Id:         "plugin-page-" + plugin.Id,
			Url:        plugin.DefaultNavUrl,
			Img:        plugin.Info.Logos.Small,
			SortWeight: dtos.WeightPlugin,
		}

		for _, include := range plugin.Includes {
			if !c.HasUserRole(include.Role) {
				continue
					}
				} else {
					link = &dtos.NavLink{
						Url:  hs.Cfg.AppSubURL + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						Text: include.Name,
					}
				}
			}

			if include.Type == "dashboard" && include.AddToNav {
				link := &dtos.NavLink{
					Url:  hs.Cfg.AppSubURL + "/dashboard/db/" + include.Slug,
					Text: include.Name,
				}
				appLink.Children = append(appLink.Children, link)
			}
		}

		if len(appLink.Children) > 0 {
			appLinks = append(appLinks, appLink)
		}
	}

	return appLinks, nil
}

func (hs *HTTPServer) getNavTree(c *models.ReqContext, hasEditPerm bool) ([]*dtos.NavLink, error) {
	hasAccess := ac.HasAccess(hs.AccessControl, c)
	navTree := []*dtos.NavLink{}

	if hasEditPerm {
		children := []*dtos.NavLink{
			{Text: "Dashboard", Icon: "apps", Url: hs.Cfg.AppSubURL + "/dashboard/new"},
		}
		if c.OrgRole == models.ROLE_ADMIN || c.OrgRole == models.ROLE_EDITOR {
			children = append(children, &dtos.NavLink{
				Text: "Folder", SubTitle: "Create a new folder to organize your dashboards", Id: "folder",
				Icon: "folder-plus", Url: hs.Cfg.AppSubURL + "/dashboards/folder/new",
			})
		}
		children = append(children, &dtos.NavLink{
			Text: "Import", SubTitle: "Import dashboard from file or Grafana.com", Id: "import", Icon: "import",
			Url: hs.Cfg.AppSubURL + "/dashboard/import",
		})
		navTree = append(navTree, &dtos.NavLink{
			Text:       "Create",
			Id:         "create",
			Icon:       "plus",
			Url:        hs.Cfg.AppSubURL + "/dashboard/new",
			Children:   children,
			SortWeight: dtos.WeightCreate,
		})
	}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "Home", Id: "home", Url: hs.Cfg.AppSubURL + "/", Icon: "home-alt", HideFromTabs: true},
		{Text: "Divider", Divider: true, Id: "divider", HideFromTabs: true},
		{Text: "Manage", Id: "manage-dashboards", Url: hs.Cfg.AppSubURL + "/dashboards", Icon: "sitemap"},
		{Text: "Playlists", Id: "playlists", Url: hs.Cfg.AppSubURL + "/playlists", Icon: "presentation-play"},
	}

	if c.IsSignedIn {
		dashboardChildNavs = append(dashboardChildNavs, &dtos.NavLink{
			Text: "Snapshots",
			Id:   "snapshots",
			Url:  hs.Cfg.AppSubURL + "/dashboard/snapshots",
			Icon: "camera",
		})

		dashboardChildNavs = append(dashboardChildNavs, &dtos.NavLink{
			Text: "Library panels",
			Id:   "library-panels",
			Url:  hs.Cfg.AppSubURL + "/library-panels",
			Icon: "library-panel",
		})
	}

	navTree = append(navTree, &dtos.NavLink{
		Id:         "dashboards",
		SubTitle:   "Manage dashboards and folders",
		Icon:       "apps",
		Url:        hs.Cfg.AppSubURL + "/",
		SortWeight: dtos.WeightDashboard,
		Children:   dashboardChildNavs,
	})

	if setting.ExploreEnabled && (c.OrgRole == models.ROLE_ADMIN || c.OrgRole == models.ROLE_EDITOR || setting.ViewersCanEdit) {
		navTree = append(navTree, &dtos.NavLink{
			Text:       "Explore",
			Id:         "explore",
			SubTitle:   "Explore your data",
			Icon:       "compass",
			SortWeight: dtos.WeightExplore,
			Url:        hs.Cfg.AppSubURL + "/explore",
		})
	}

	if c.IsSignedIn {
		navTree = append(navTree, hs.getProfileNode(c))
	}

	if setting.AlertingEnabled {
		alertChildNavs := []*dtos.NavLink{
			{Text: "Alert rules", Id: "alert-list", Url: hs.Cfg.AppSubURL + "/alerting/list", Icon: "list-ul"},
		}
		if hs.Cfg.IsNgAlertEnabled() {
			alertChildNavs = append(alertChildNavs, &dtos.NavLink{Text: "Silences", Id: "silences", Url: hs.Cfg.AppSubURL + "/alerting/silences", Icon: "bell-slash"})
		}
		if c.OrgRole == models.ROLE_ADMIN || c.OrgRole == models.ROLE_EDITOR {
			if hs.Cfg.IsNgAlertEnabled() {
				alertChildNavs = append(alertChildNavs, &dtos.NavLink{
					Text: "Contact points", Id: "receivers", Url: hs.Cfg.AppSubURL + "/alerting/notifications",
					Icon: "comment-alt-share",
				})
				alertChildNavs = append(alertChildNavs, &dtos.NavLink{Text: "Notification policies", Id: "am-routes", Url: hs.Cfg.AppSubURL + "/alerting/routes", Icon: "sitemap"})
			} else {
				alertChildNavs = append(alertChildNavs, &dtos.NavLink{
					Text: "Notification channels", Id: "channels", Url: hs.Cfg.AppSubURL + "/alerting/notifications",
					Icon: "comment-alt-share",
				})
			}
		}

		navTree = append(navTree, &dtos.NavLink{
			Text:       "Alerting",
			SubTitle:   "Alert rules and notifications",
			Id:         "alerting",
			Icon:       "bell",
			Url:        hs.Cfg.AppSubURL + "/alerting/list",
			Children:   alertChildNavs,
			SortWeight: dtos.WeightAlerting,
		})
	}

	appLinks, err := hs.getAppLinks(c)

	configNodes := []*dtos.NavLink{}

	if c.OrgRole == models.ROLE_ADMIN {
		configNodes = append(configNodes, &dtos.NavLink{
			Text:        "Data sources",
			Icon:        "database",
		})
	}

	if hasAccess(ac.ReqOrgAdmin, ac.ActionOrgUsersRead, ac.ScopeUsersAll) {
		configNodes = append(configNodes, &dtos.NavLink{
			Text:        "Users",
			Id:          "users",
		})
	}

	if c.OrgRole == models.ROLE_ADMIN || (hs.Cfg.EditorsCanAdmin && c.OrgRole == models.ROLE_EDITOR) {
		configNodes = append(configNodes, &dtos.NavLink{
			Text:        "Teams",
			Id:          "teams",
			Icon:        "plug",
			Url:         hs.Cfg.AppSubURL + "/plugins",
		})

		configNodes = append(configNodes, &dtos.NavLink{
			Text:        "Preferences",
			Id:          "org-settings",
			Icon:        "sliders-v-alt",
			Url:         hs.Cfg.AppSubURL + "/org",
		})
		configNodes = append(configNodes, &dtos.NavLink{
			Text:        "API keys",
			Id:          "apikeys",
			Url:         hs.Cfg.AppSubURL + "/org/apikeys",
		})
	}

	if len(configNodes) > 0 {
		navTree = append(navTree, &dtos.NavLink{
			Id:         dtos.NavIDCfg,
			Text:       "Configuration",
			SubTitle:   "Organization: " + c.OrgName,
			Url:        configNodes[0].Url,
			SortWeight: dtos.WeightConfig,
			Children:   configNodes,
		})
	}

	adminNavLinks := hs.buildAdminNavLinks(c)

	if len(adminNavLinks) > 0 {
		navTree = append(navTree, &dtos.NavLink{
			Text:         "Server Admin",
			SubTitle:     "Manage all users and orgs",
			HideFromTabs: true,
			Id:           "admin",
			Icon:         "shield",
			Url:          adminNavLinks[0].Url,
			SortWeight:   dtos.WeightAdmin,
			Children:     adminNavLinks,
		})
	}

	helpVersion := fmt.Sprintf(`%s v%s (%s)`, setting.ApplicationName, setting.BuildVersion, setting.BuildCommit)
	if hs.Cfg.AnonymousHideVersion && !c.IsSignedIn {
		helpVersion = setting.ApplicationName
	}

	navTree = append(navTree, &dtos.NavLink{
		Text:         "Help",
		SubTitle:     helpVersion,
		Id:           "help",
		Url:          "#",
		Icon:         "question-circle",
		HideFromMenu: true,
		SortWeight:   dtos.WeightHelp,
		Children:     []*dtos.NavLink{},
	})

	return navTree, nil
}

func (hs *HTTPServer) buildAdminNavLinks(c *models.ReqContext) []*dtos.NavLink {
	hasAccess := ac.HasAccess(hs.AccessControl, c)
	adminNavLinks := []*dtos.NavLink{}

	if hasAccess(ac.ReqGrafanaAdmin, ac.ActionUsersRead, ac.ScopeGlobalUsersAll) {
		adminNavLinks = append(adminNavLinks, &dtos.NavLink{
			Text: "Users", Id: "global-users", Url: hs.Cfg.AppSubURL + "/admin/users", Icon: "user",
		})
	}

	if c.IsGrafanaAdmin {
		adminNavLinks = append(adminNavLinks, &dtos.NavLink{
			Text: "Orgs", Id: "global-orgs", Url: hs.Cfg.AppSubURL + "/admin/orgs", Icon: "building",
		})
		adminNavLinks = append(adminNavLinks, &dtos.NavLink{
			Text: "Settings", Id: "server-settings", Url: hs.Cfg.AppSubURL + "/admin/settings", Icon: "sliders-v-alt",
		})
		adminNavLinks = append(adminNavLinks, &dtos.NavLink{
			Text: "Stats", Id: "server-stats", Url: hs.Cfg.AppSubURL + "/admin/stats", Icon: "graph-bar",
		})
	}

	if hs.Cfg.LDAPEnabled && hasAccess(ac.ReqGrafanaAdmin, ac.ActionLDAPStatusRead) {
		adminNavLinks = append(adminNavLinks, &dtos.NavLink{
			Text: "LDAP", Id: "ldap", Url: hs.Cfg.AppSubURL + "/admin/ldap", Icon: "book",
		})
	}

	return adminNavLinks
}

func (hs *HTTPServer) setIndexViewData(c *models.ReqContext) (*dtos.IndexViewData, error) {
	hasEditPermissionInFoldersQuery := models.HasEditPermissionInFoldersQuery{SignedInUser: c.SignedInUser}
	if err := bus.Dispatch(&hasEditPermissionInFoldersQuery); err != nil {
		return nil, err
	}
	hasEditPerm := hasEditPermissionInFoldersQuery.Result

	settings, err := hs.getFrontendSettingsMap(c)
	if err != nil {
	settings["dateFormats"] = hs.Cfg.DateFormats

	prefsQuery := models.GetPreferencesWithDefaultsQuery{User: c.SignedInUser}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			OrgCount:                   c.OrgCount,
			OrgId:                      c.OrgId,
			IsGrafanaAdmin:             c.IsGrafanaAdmin,
			LightTheme:                 prefs.Theme == lightName,
			Timezone:                   prefs.Timezone,
			Locale:                     locale,
			HelpFlags1:                 c.HelpFlags1,
			HasEditPermissionInFolders: hasEditPerm,
		GoogleTagManagerId:      setting.GoogleTagManagerId,
		BuildVersion:            setting.BuildVersion,
		BuildCommit:             setting.BuildCommit,
		NewGrafanaVersion:       hs.PluginManager.GrafanaLatestVersion(),
		NewGrafanaVersionExists: hs.PluginManager.GrafanaHasUpdate(),
		AppName:                 setting.ApplicationName,
		AppNameBodyClass:        getAppNameBodyClass(hs.License.HasValidLicense()),
		FavIcon:                 "public/img/fav32.png",
		AppleTouchIcon:          "public/img/apple-touch-icon.png",
		AppTitle:                "Grafana",
		Sentry:                  &hs.Cfg.Sentry,
		Nonce:                   c.RequestNonce,
		ContentDeliveryURL:      hs.Cfg.GetContentDeliveryURL(hs.License.ContentDeliveryPrefix()),
	}

	if hs.Cfg.FeatureToggles["accesscontrol"] {
		userPermissions, err := hs.AccessControl.GetUserPermissions(c.Req.Context(), c.SignedInUser)
		if err != nil {
			return nil, err
		}

	c.HTML(404, "index", data)
}

func getAppNameBodyClass(validLicense bool) string {
	if validLicense {
		return "app-enterprise"
	}

	return "app-grafana"
}
