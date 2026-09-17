package main

package packaging

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grafana/grafana/pkg/build/config"
	"github.com/grafana/grafana/pkg/build/errutil"
	"github.com/grafana/grafana/pkg/build/grafana"
	"github.com/grafana/grafana/pkg/build/plugins"
	"github.com/grafana/grafana/pkg/build/syncutil"
	"github.com/grafana/grafana/pkg/infra/fs"
)

var (
	ErrorNoBinaries = errors.New("no binaries found")
	ErrorNoDebArch  = errors.New("deb architecture not defined")
	ErrorNoRPMArch  = errors.New("rpm architecture not defined")
)

const (
	maxAttempts          = 3
	enterpriseSfx        = "-enterprise"
	enterprise2Sfx       = "-enterprise2"
	DefaultDebDBBucket   = "grafana-aptly-db"
	DefaultDebRepoBucket = "grafana-repo"
	DefaultRPMRepoBucket = "grafana-repo"
	DefaultTTLSeconds    = "300"
)

// PackageRegexp returns a regexp for matching packages corresponding to a certain Grafana edition.
func PackageRegexp(edition config.Edition) *regexp.Regexp {
	var sfx string
	switch edition {
	case config.EditionOSS:
	case config.EditionEnterprise:
		sfx = "-enterprise"
	case config.EditionEnterprise2:
		sfx = "-enterprise2"
	default:
		panic(fmt.Sprintf("unrecognized edition %q", edition))
	}
	rePkg, err := regexp.Compile(fmt.Sprintf(`^grafana%s(?:-rpi)?[-_][^-_]+.*$`, sfx))
	if err != nil {
		panic(fmt.Sprintf("Failed to compile regexp: %s", err))
	}

	return rePkg
}

// PackageGrafana packages Grafana for various variants.
func PackageGrafana(
	ctx context.Context,
	version string,
	grafanaDir string,
	cfg config.Config,
	edition config.Edition,
	variants []config.Variant,
	shouldSign bool,
	p syncutil.WorkerPool,
) error {
	if err := packageGrafana(ctx, edition, version, grafanaDir, variants, shouldSign, p); err != nil {
		return err
	}
	if err := signRPMPackages(edition, cfg, grafanaDir); err != nil {
		return err
	}
	if err := checksumPackages(grafanaDir, edition); err != nil {
		return err
	}

	return nil
}

func packageGrafana(
	ctx context.Context,
	edition config.Edition,
	version string,
	grafanaDir string,
	variants []config.Variant,
	shouldSign bool,
	p syncutil.WorkerPool,
) error {
	distDir := filepath.Join(grafanaDir, "dist")
	exists, err := fs.Exists(distDir)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf("directory %s doesn't exist - creating...", distDir)
		//nolint
		if err := os.MkdirAll(distDir, 0o755); err != nil {
			return fmt.Errorf("couldn't create dist: %w", err)
		}
	}

	var m pluginsManifest
	manifestPath := filepath.Join(grafanaDir, "plugins-bundled", "external.json")
	//nolint:gosec
	manifestB, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to open plugins manifest %q: %w", manifestPath, err)
	}
	if err := json.Unmarshal(manifestB, &m); err != nil {
		return err
	}

	g, ctx := errutil.GroupWithContext(ctx)
	for _, v := range variants {
		packageVariant(ctx, v, edition, version, grafanaDir, shouldSign, g, p, m)
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// packageVariant packages Grafana for a certain variant.
func packageVariant(
	ctx context.Context,
	v config.Variant,
	edition config.Edition,
	version string,
	grafanaDir string,
	shouldSign bool,
	g *errutil.Group,
	p syncutil.WorkerPool,
	m pluginsManifest,
) {
	p.Schedule(g.Wrap(func() error {
		// We've experienced spurious packaging failures, so retry on failure.
		i := 0
		for {
			if err := realPackageVariant(ctx, v, edition, version, grafanaDir, m, shouldSign); err != nil {
				i++
				if i < maxAttempts {
					log.Printf("Packaging for variant %s, %s edition failed: %s, trying again", v, edition, err)
					continue
				}

				return err
			}

			break
		}

		return nil
	}))
}

// signRPMPackages signs the RPM packages.
func signRPMPackages(edition config.Edition, cfg config.Config, grafanaDir string) error {
	log.Printf("Signing %s RPM packages...", edition)
	var sfx string
	switch edition {
	case config.EditionOSS:
	case config.EditionEnterprise:
		sfx = enterpriseSfx
	case config.EditionEnterprise2:
		sfx = enterprise2Sfx
	default:
		panic(fmt.Sprintf("Unrecognized edition %s", edition))
	}
	rpms, err := filepath.Glob(filepath.Join(grafanaDir, "dist", fmt.Sprintf("grafana%s-*.rpm", sfx)))
	if err != nil {
		return err
	}

	rpmArgs := append([]string{"--addsign"}, rpms...)
	log.Printf("Invoking rpm with args: %+v", rpmArgs)
	//nolint:gosec
	cmd := exec.Command("rpm", rpmArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to sign RPM packages: %s", output)
	}
	if err := os.Remove(cfg.GPGPassPath); err != nil {
		return fmt.Errorf("failed to remove %q: %w", cfg.GPGPassPath, err)
	}

	log.Printf("Verifying %s RPM packages...", edition)
	// The output changed between rpm versions
	reOutput := regexp.MustCompile("(?:digests signatures OK)|(?:pgp.+OK)")
	for _, p := range rpms {
		//nolint:gosec
		cmd := exec.Command("rpm", "-K", p)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to verify RPM signature: %w", err)
		}

		if !reOutput.Match(output) {
			return fmt.Errorf("RPM package %q not verified: %s", p, output)
		}
	}

	return nil
}

// checksumPackages generates package checksums with SHA-256.
func checksumPackages(grafanaDir string, edition config.Edition) error {
	log.Printf("Checksumming %s packages...", edition)
	distDir := filepath.Join(grafanaDir, "dist")
	rePkg := PackageRegexp(edition)
	if err := filepath.Walk(distDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			var pathErr *os.PathError
			if errors.As(err, &pathErr) {
				log.Printf("path error in walk function for file %q: %s", pathErr.Path, pathErr.Err.Error())
				return nil
			}
			return fmt.Errorf("walking through dist folder failed: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		fname := filepath.Base(fpath)
		if strings.HasSuffix(fname, ".sha256") || strings.HasSuffix(fname, ".version") || !rePkg.MatchString(fname) {
			log.Printf("Ignoring non-package %q", fpath)
			return nil
		}

		return shaFile(fpath)
	}); err != nil {
		return fmt.Errorf("checksumming packages in %q failed: %w", distDir, err)
	}

	log.Printf("Successfully checksummed %s packages", edition)
	return nil
}

func shaFile(fpath string) error {
	//nolint:gosec
	fd, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("failed to open %q: %w", fpath, err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			log.Println(err)
		}
	}()

	h := sha256.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", fpath, err)
	}

	//nolint:gosec
	out, err := os.Create(fpath + ".sha256")
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", fpath+".sha256", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println("failed to close file", out.Name())
		}
	}()

	if _, err = fmt.Fprintf(out, "%x\n", h.Sum(nil)); err != nil {
		return fmt.Errorf("failed to write %q: %w", out.Name(), err)
	}

	return nil
}

// createPackage creates a Linux package.
func createPackage(srcDir string, options linuxPackageOptions) error {
	cliBinary := "grafana-cli"
	serverBinary := "grafana-server"

	packageRoot, err := os.MkdirTemp("", "grafana-linux-pack")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(packageRoot); err != nil {
			log.Println(err)
		}
	}()

	for _, dname := range []string{
		options.homeDir,
		options.configDir,
		"etc/init.d",
		options.etcDefaultPath,
		"usr/lib/systemd/system",
		"usr/sbin",
	} {
		dpath := filepath.Join(packageRoot, dname)
		//nolint
		if err := os.MkdirAll(dpath, 0o755); err != nil {
			return fmt.Errorf("failed to make directory %q: %w", dpath, err)
		}
	}

	if err := fs.CopyFile(options.cliBinaryWrapperSrc, filepath.Join(packageRoot, "usr", "sbin", cliBinary)); err != nil {
		return err
	}
	if err := fs.CopyFile(filepath.Join(srcDir, "bin", serverBinary),
		filepath.Join(packageRoot, "usr", "sbin", serverBinary)); err != nil {
		return err
	}
	if err := fs.CopyFile(options.initdScriptSrc, filepath.Join(packageRoot, options.initdScriptFilePath)); err != nil {
		return err
	}
	if err := fs.CopyFile(options.defaultFileSrc, filepath.Join(packageRoot, options.etcDefaultFilePath)); err != nil {
		return err
	}
	if err := fs.CopyFile(options.systemdFileSrc, filepath.Join(packageRoot, options.systemdServiceFilePath)); err != nil {
		return err
	}
	if err := fs.CopyRecursive(srcDir, filepath.Join(packageRoot, options.homeDir)); err != nil {
		return err
	}
	homeBinDir := filepath.Join(packageRoot, options.homeBinDir)
	if err := os.RemoveAll(homeBinDir); err != nil {
		return fmt.Errorf("failed to remove %q: %w", homeBinDir, err)
	}
	//nolint
	if err := os.MkdirAll(homeBinDir, 0o755); err != nil {
		return fmt.Errorf("failed to make directory %q: %w", homeBinDir, err)
	}
	// The grafana-cli binary is exposed through a wrapper to ensure a proper
	// configuration is in place. To enable that, we need to store the original
	// binary in a separate location to avoid conflicts.
	if err := fs.CopyFile(filepath.Join(srcDir, "bin", cliBinary), filepath.Join(homeBinDir, cliBinary)); err != nil {
		return err
	}

	if err := executeFPM(options, packageRoot, srcDir); err != nil {
		return err
	}

	return nil
}
func executeFPM(options linuxPackageOptions, packageRoot, srcDir string) error {
	name := "grafana"
	vendor := "Grafana"
	if options.edition == config.EditionEnterprise || options.edition == config.EditionEnterprise2 {
		vendor += " Enterprise"
		if options.edition == config.EditionEnterprise2 {
			name += enterprise2Sfx
		} else if options.edition == config.EditionEnterprise {
			name += enterpriseSfx
		}
	}

	if options.goArch == config.ArchARM && options.goArm == "6" {
		name += "-rpi"
	}

	pkgVersion := packageVersion(options)
	args := []string{
		"-s", "dir",
		"--description", "Grafana",
		"-C", packageRoot,
		"--url", "https://grafana.com",
		"--maintainer", "contact@grafana.com",
		"--config-files", options.initdScriptFilePath,
		"--config-files", options.etcDefaultFilePath,
		"--config-files", options.systemdServiceFilePath,
		"--after-install", options.postinstSrc,
		"--version", pkgVersion,
		"-p", "dist/",
		"--name", name,
		"--vendor", vendor,
		"-a", string(options.packageArch),
	}
	if options.edition == config.EditionEnterprise || options.edition == config.EditionEnterprise2 || options.goArch == config.ArchARMv6 {
		args = append(args, "--replaces", "grafana")
	}
	if options.edition == config.EditionOSS {
		args = append(args, "--license", "\"AGPLv3\"")
	}
	switch options.packageType {
	case packageTypeRpm:
		args = append(args, "-t", "rpm", "--rpm-posttrans", "packaging/rpm/control/posttrans")
	case packageTypeDeb:
		args = append(args, "-t", "deb", "--deb-no-default-config-files")
	default:
		panic(fmt.Sprintf("Unrecognized package type %d", options.packageType))
	}
	for _, dep := range options.depends {
		args = append(args, "--depends", dep)
	}
	args = append(args, ".")

	distDir := filepath.Join(options.grafanaDir, "dist")
	log.Printf("Generating package in %q (source directory %q)", distDir, srcDir)

	cmdStr := "fpm"
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			arg = fmt.Sprintf("'%s'", arg)
		}
		cmdStr += fmt.Sprintf(" %s", arg)
	}
	log.Printf("Creating %s package: %s...", options.packageType, cmdStr)
	const rvmPath = "/etc/profile.d/rvm.sh"
	exists, err := fs.Exists(rvmPath)
	if err != nil {
		return err
	}
	if exists {
		cmdStr = fmt.Sprintf("source %q && %s", rvmPath, cmdStr)
		log.Printf("Sourcing %q before running fpm", rvmPath)
	}
	//nolint:gosec
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Dir = options.grafanaDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run fpm: %s", output)
	}

	return nil
}

// copyPubDir copies public/ from grafanaDir to tmpDir.
func copyPubDir(grafanaDir, tmpDir string) error {
	srcPubDir := filepath.Join(grafanaDir, "public")
	tgtPubDir := filepath.Join(tmpDir, "public")
	log.Printf("Copying %q to %q...", srcPubDir, tgtPubDir)
	if err := fs.CopyRecursive(srcPubDir, tgtPubDir); err != nil {
		return fmt.Errorf("failed to copy %q to %q: %w", srcPubDir, tgtPubDir, err)
	}

	return nil
}

// copyBinaries copies binaries from grafanaDir into tmpDir.
func copyBinaries(grafanaDir, tmpDir string, args grafana.BuildArgs, edition config.Edition) error {
	tgtDir := filepath.Join(tmpDir, "bin")
	//nolint
	if err := os.MkdirAll(tgtDir, 0o755); err != nil {
		return fmt.Errorf("failed to make directory %q: %w", tgtDir, err)
	}

	binDir := filepath.Join(grafanaDir, "bin", grafana.BinaryFolder(edition, args))

	files, err := os.ReadDir(binDir)
	if err != nil {
		return fmt.Errorf("failed to list files in %q: %w", binDir, err)
	}

	if len(files) == 0 {
		return fmt.Errorf("%w in %s", ErrorNoBinaries, binDir)
	}

	for _, file := range files {
		srcPath := filepath.Join(binDir, file.Name())
		tgtPath := filepath.Join(tgtDir, file.Name())

		if err := fs.CopyFile(srcPath, tgtPath); err != nil {
			return err
		}
	}

	return nil
}

// copyScripts copies scripts from grafanaDir into tmpDir.
func copyScripts(grafanaDir, tmpDir string) error {
	//nolint
	if err := os.MkdirAll(filepath.Join(tmpDir, "scripts"), 0o755); err != nil {
		return fmt.Errorf("failed to create dir %q: %w", filepath.Join(tmpDir, "scripts"), err)
	}
	scriptsDir := filepath.Join(grafanaDir, "scripts")
	infos, err := os.ReadDir(scriptsDir)
	if err != nil {
		return fmt.Errorf("failed to list files in %q: %w", scriptsDir, err)
	}
	for _, file := range infos {
		info, err := file.Info()
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}

		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}

		path := ""

		path = filepath.Join(scriptsDir, info.Name())

		if err := fs.CopyFile(path, filepath.Join(tmpDir, "scripts", info.Name())); err != nil {
			return fmt.Errorf("failed to copy %q to %q: %w", path, tmpDir, err)
		}
	}

	return nil
}

// copyConfFiles copies configuration files from grafanaDir into tmpDir.
func copyConfFiles(grafanaDir, tmpDir string) error {
	//nolint:gosec
	if err := os.MkdirAll(filepath.Join(tmpDir, "conf"), 0o755); err != nil {
		return fmt.Errorf("failed to create dir %q: %w", filepath.Join(tmpDir, "conf"), err)
	}

	confDir := filepath.Join(grafanaDir, "conf")
	infos, err := os.ReadDir(confDir)
	if err != nil {
		return fmt.Errorf("failed to list files in %q: %w", confDir, err)
	}
	for _, info := range infos {
		fpath := filepath.Join(confDir, info.Name())
		if err := fs.CopyRecursive(fpath, filepath.Join(tmpDir, "conf", info.Name())); err != nil {
			return err
		}
	}

	return nil
}

// copyPlugins copies plugins from grafanaDir into tmpDir.
func copyPlugins(ctx context.Context, v config.Variant, grafanaDir, tmpDir string, m pluginsManifest, shouldSign bool) error {
	log.Printf("Copying plugins for package variant %s...", v)

	variant2Sfx := map[config.Variant]string{
		config.VariantLinuxAmd64:   "linux_amd64",
		config.VariantDarwinAmd64:  "darwin_amd64",
		config.VariantWindowsAmd64: "windows_amd64.exe",
	}

	tgtDir := filepath.Join(tmpDir, "plugins-bundled")
	exists, err := fs.Exists(tgtDir)
	if err != nil {
		return err
	}
	if !exists {
		//nolint:gosec
		if err := os.MkdirAll(tgtDir, 0o755); err != nil {
			return err
		}
	}
	pluginsDir := filepath.Join(grafanaDir, "plugins-bundled")

	// External plugins.
	for _, pm := range m.Plugins {
		srcDir := filepath.Join(pluginsDir, fmt.Sprintf("%s-%s", pm.Name, pm.Version))
		dstDir := filepath.Join(tgtDir, fmt.Sprintf("%s-%s", pm.Name, pm.Version))
		log.Printf("Copying external plugin %q to %q...", srcDir, dstDir)

		//nolint:gosec
		jsonB, err := os.ReadFile(filepath.Join(srcDir, "plugin.json"))
		if err != nil {
			return fmt.Errorf("failed to read %q: %w", filepath.Join(srcDir, "plugin.json"), err)
		}
		var plugJSON map[string]interface{}
		if err := json.Unmarshal(jsonB, &plugJSON); err != nil {
			return err
		}

		plugExe, ok := plugJSON["executable"].(string)
		var wantExe string
		if ok && strings.TrimSpace(plugExe) != "" {
			sfx := variant2Sfx[v]
			if sfx == "" {
				log.Printf("External plugin %s-%s doesn't have an executable for variant %s - ignoring",
					pm.Name, pm.Version, v)
				continue
			}

			wantExe = fmt.Sprintf("%s_%s", plugExe, sfx)
			log.Printf("The external plugin should contain an executable %q", wantExe)
			exists, err := fs.Exists(filepath.Join(srcDir, wantExe))
			if err != nil {
				return err
			}
			if !exists {
				log.Printf("External plugin %s-%s doesn't have an executable of the right format: %q - ignoring",
					pm.Name, pm.Version, wantExe)
				continue
			}
		}

		if err := filepath.Walk(srcDir, func(pth string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			log.Printf("Handling %q", pth)

			relPath := strings.TrimPrefix(pth, srcDir)
			relPath = strings.TrimPrefix(relPath, "/")
			dstPath := filepath.Join(dstDir, relPath)

			if info.IsDir() {
				log.Printf("Making directory %q", dstPath)
				//nolint:gosec
				return os.MkdirAll(dstPath, info.Mode())
			}

			if wantExe != "" {
				m, err := regexp.MatchString(fmt.Sprintf(`^%s_[^/]+$`, plugExe), relPath)
				if err != nil {
					return err
				}
				if m && relPath != wantExe {
					// Ignore other executable variants
					log.Printf("Ignoring executable variant %q", pth)
					return nil
				}
			}

			log.Printf("Copying %q to %q", pth, dstPath)
			return fs.CopyFile(pth, dstPath)
		}); err != nil {
			return fmt.Errorf("failed to copy external plugin %q to %q: %w", srcDir, dstDir, err)
		}

		if shouldSign {
			if err := plugins.BuildManifest(ctx, dstDir, true); err != nil {
				return fmt.Errorf("failed to generate signed manifest for external plugin %q: %w", dstDir, err)
			}
		}
	}

	return copyInternalPlugins(pluginsDir, tmpDir)
}

func copyInternalPlugins(pluginsDir, tmpDir string) error {
	tgtDir := filepath.Join(tmpDir, "plugins-bundled", "internal")
	srcDir := filepath.Join(pluginsDir, "dist")

	exists, err := fs.Exists(tgtDir)
	if err != nil {
		return err
	}
	if !exists {
		//nolint:gosec
		if err := os.MkdirAll(tgtDir, 0o755); err != nil {
			return err
		}
	}

	// Copy over internal plugins.
	fis, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to list internal plugins in %q: %w", srcDir, err)
	}
	for _, fi := range fis {
		srcPath := filepath.Join(srcDir, fi.Name())
		if !fi.IsDir() {
			log.Printf("Ignoring non-directory %q", srcPath)
			continue
		}

		dstPath := filepath.Join(tgtDir, fi.Name())
		log.Printf("Copying internal plugin %q to %q...", srcPath, dstPath)
		if err := fs.CopyRecursive(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy %q to %q: %w", srcPath, dstPath, err)
		}
	}

	return nil
}

func realPackageVariant(ctx context.Context, v config.Variant, edition config.Edition, version, grafanaDir string,
	m pluginsManifest, shouldSign bool) error {
	log.Printf("Packaging Grafana %s for %s...", edition, v)

	enableDeb := false
	enableRpm := false
	switch v {
	case config.VariantLinuxAmd64:
		enableDeb = true
		enableRpm = true
	case config.VariantArmV6:
		enableDeb = true
	case config.VariantArmV7:
		enableDeb = true
		enableRpm = true
	case config.VariantArm64:
		enableDeb = true
		enableRpm = true
	default:
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Println(err)
		}
	}()

	args := grafana.VariantBuildArgs(v)

	if err := copyPubDir(grafanaDir, tmpDir); err != nil {
		return err
	}
	if err := copyBinaries(grafanaDir, tmpDir, args, edition); err != nil {
		return err
	}
	if err := copyScripts(grafanaDir, tmpDir); err != nil {
		return err
	}
	if err := copyConfFiles(grafanaDir, tmpDir); err != nil {
		return err
	}
	if err := copyPlugins(ctx, v, grafanaDir, tmpDir, m, shouldSign); err != nil {
		return err
	}

	if v == config.VariantWindowsAmd64 {
		toolsDir := filepath.Join(tmpDir, "tools")
		//nolint:gosec
		if err := os.MkdirAll(toolsDir, 0o755); err != nil {
			return fmt.Errorf("failed to create tools dir %q: %w", toolsDir, err)
		}

		if err := fs.CopyFile("/usr/local/go/lib/time/zoneinfo.zip",
			filepath.Join(tmpDir, "tools", "zoneinfo.zip")); err != nil {
			return err
		}
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "VERSION"), []byte(version), 0664); err != nil {
		return fmt.Errorf("failed to write %s/VERSION: %w", tmpDir, err)
	}

	if err := createArchive(tmpDir, edition, v, version, grafanaDir); err != nil {
		return err
	}

	if enableDeb {
		if args.DebArch == "" {
			return fmt.Errorf("%w for %s", ErrorNoDebArch, v)
		}

		if err := createPackage(tmpDir, linuxPackageOptions{
			edition:                edition,
			version:                version,
			grafanaDir:             grafanaDir,
			goArch:                 args.GoArch,
			goArm:                  args.GoArm,
			packageType:            packageTypeDeb,
			packageArch:            args.DebArch,
			homeDir:                "/usr/share/grafana",
			homeBinDir:             "/usr/share/grafana/bin",
			binPath:                "/usr/sbin",
			configDir:              "/etc/grafana",
			etcDefaultPath:         "/etc/default",
			etcDefaultFilePath:     "/etc/default/grafana-server",
			initdScriptFilePath:    "/etc/init.d/grafana-server",
			systemdServiceFilePath: "/usr/lib/systemd/system/grafana-server.service",
			postinstSrc:            filepath.Join(grafanaDir, "packaging", "deb", "control", "postinst"),
			initdScriptSrc:         filepath.Join(grafanaDir, "packaging", "deb", "init.d", "grafana-server"),
			defaultFileSrc:         filepath.Join(grafanaDir, "packaging", "deb", "default", "grafana-server"),
			systemdFileSrc:         filepath.Join(grafanaDir, "packaging", "deb", "systemd", "grafana-server.service"),
			cliBinaryWrapperSrc:    filepath.Join(grafanaDir, "packaging", "wrappers", "grafana-cli"),
			depends:                []string{"adduser", "libfontconfig1"},
		}); err != nil {
			return err
		}
	}

	if !enableRpm {
		return nil
	}

	if args.RPMArch == "" {
		return fmt.Errorf("%w for %s", ErrorNoRPMArch, v)
	}

	if err := createPackage(tmpDir, linuxPackageOptions{
		edition:                edition,
		version:                version,
		grafanaDir:             grafanaDir,
		goArch:                 args.GoArch,
		packageType:            packageTypeRpm,
		packageArch:            args.RPMArch,
		homeDir:                "/usr/share/grafana",
		homeBinDir:             "/usr/share/grafana/bin",
		binPath:                "/usr/sbin",
		configDir:              "/etc/grafana",
		etcDefaultPath:         "/etc/sysconfig",
		etcDefaultFilePath:     "/etc/sysconfig/grafana-server",
		initdScriptFilePath:    "/etc/init.d/grafana-server",
		systemdServiceFilePath: "/usr/lib/systemd/system/grafana-server.service",
		postinstSrc:            filepath.Join(grafanaDir, "packaging", "rpm", "control", "postinst"),
		initdScriptSrc:         filepath.Join(grafanaDir, "packaging", "rpm", "init.d", "grafana-server"),
		defaultFileSrc:         filepath.Join(grafanaDir, "packaging", "rpm", "sysconfig", "grafana-server"),
		systemdFileSrc:         filepath.Join(grafanaDir, "packaging", "rpm", "systemd", "grafana-server.service"),
		cliBinaryWrapperSrc:    filepath.Join(grafanaDir, "packaging", "wrappers", "grafana-cli"),
		// chkconfig is depended on since our systemd service wraps a SysV init script, and that requires chkconfig
		depends: []string{"/sbin/service", "chkconfig", "fontconfig", "freetype", "urw-fonts"},
	}); err != nil {
		return err
	}

	return nil
}

// pluginManifest has details of an external plugin package.
type pluginManifest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}

// pluginsManifest represents a manifest of Grafana's external plugins.
type pluginsManifest struct {
	Plugins []pluginManifest `json:"plugins"`
}

// packageVersion converts a Grafana version into the corresponding package version.
func packageVersion(options linuxPackageOptions) string {
	verComponents := strings.Split(options.version, "-")
	if len(verComponents) > 2 {
		panic(fmt.Sprintf("Version string contains more than one hyphen: %q", options.version))
	}

	switch options.packageType {
	case packageTypeDeb, packageTypeRpm:
		if len(verComponents) > 1 {
			// With Debian and RPM packages, it's customary to prefix any pre-release component with a ~, since this
			// is considered of lower lexical value than the empty character, and this way pre-release versions are
			// considered to be of a lower version than the final version (which lacks this suffix).
			return fmt.Sprintf("%s~%s", verComponents[0], verComponents[1])
		}

		return options.version
	default:
		panic(fmt.Sprintf("Unrecognized package type %s", options.packageType))
	}
}

type packageType int

func (pt packageType) String() string {
	switch pt {
	case packageTypeDeb:
		return "Debian"
	case packageTypeRpm:
		return "RPM"
	default:
		panic(fmt.Sprintf("Unrecognized package type %d", pt))
	}
}

const (
	packageTypeDeb packageType = iota
	packageTypeRpm
)

type linuxPackageOptions struct {
	edition                config.Edition
	packageType            packageType
	version                string
	grafanaDir             string
	goArch                 config.Architecture
	goArm                  string
	packageArch            config.Architecture
	homeDir                string
	homeBinDir             string
	binPath                string
	configDir              string
	etcDefaultPath         string
	etcDefaultFilePath     string
	initdScriptFilePath    string
	systemdServiceFilePath string
	postinstSrc            string
	initdScriptSrc         string
	defaultFileSrc         string
	systemdFileSrc         string
	cliBinaryWrapperSrc    string

	depends []string
}

// createArchive makes a distribution archive.
func createArchive(srcDir string, edition config.Edition, v config.Variant, version, grafanaDir string) error {
	distDir := filepath.Join(grafanaDir, "dist")
	exists, err := fs.Exists(distDir)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf("directory %s doesn't exist - creating...", distDir)
		//nolint:gosec
		if err := os.MkdirAll(distDir, 0o755); err != nil {
			return fmt.Errorf("couldn't create dist: %w", err)
		}
	}
	sfx := ""
	if edition == config.EditionEnterprise2 {
		sfx = enterprise2Sfx
	} else if edition == config.EditionEnterprise {
		sfx = enterpriseSfx
	}
	if v != config.VariantWindowsAmd64 {
		return createTarball(srcDir, version, string(v), sfx, grafanaDir)
	}

	return createZip(srcDir, version, string(v), sfx, grafanaDir)
}

func createZip(srcDir, version, variantStr, sfx, grafanaDir string) error {
	fpath := filepath.Join(grafanaDir, "dist", fmt.Sprintf("grafana%s-%s.%s.zip", sfx, version, variantStr))
	//nolint:gosec
	tgt, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", fpath, err)
	}
	defer func() {
		if err := tgt.Close(); err != nil {
			log.Println(err)
		}
	}()

	//nolint:gosec
	if err := os.Chmod(fpath, 0664); err != nil {
		return fmt.Errorf("failed to set permissions on %q: %w", fpath, err)
	}
	zipWriter := zip.NewWriter(tgt)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, fname := range []string{"LICENSE", "README.md", "NOTICE.md"} {
		fpath := filepath.Join(grafanaDir, fname)
		fi, err := os.Lstat(fpath)
		if err != nil {
			return fmt.Errorf("couldn't stat %q: %w", fpath, err)
		}
		hdr, err := zip.FileInfoHeader(fi)
		if err != nil {
			return fmt.Errorf("failed to open zip header: %w", err)
		}
		// Enable compression, as it's disabled by default
		hdr.Method = zip.Deflate
		hdr.Name = fmt.Sprintf("grafana-%s/%s", version, fname)
		w, err := zipWriter.CreateHeader(hdr)
		if err != nil {
			return fmt.Errorf("failed writing zip header: %w", err)
		}
		//nolint:gosec
		src, err := os.Open(fpath)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", fname, err)
		}
		if _, err := io.Copy(w, src); err != nil {
			if err := src.Close(); err != nil {
				log.Println(err)
			}
			return fmt.Errorf("failed writing zip entry: %w", err)
		}
		if err := src.Close(); err != nil {
			log.Println(err)
		}
	}
	if err := filepath.Walk(srcDir, func(fpath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fpath == srcDir {
			return nil
		}

		hdr, err := zip.FileInfoHeader(fi)
		if err != nil {
			return fmt.Errorf("failed to open zip header: %s", err)
		}
		// Enable compression, as it's disabled by default
		hdr.Method = zip.Deflate
		hdr.Name = fmt.Sprintf("grafana-%s/%s", version, strings.TrimPrefix(fpath, fmt.Sprintf("%s/", srcDir)))
		if fi.IsDir() {
			// A trailing slash means it's a directory
			if hdr.Name[len(hdr.Name)-1] != '/' {
				hdr.Name += "/"
			}
		}
		w, err := zipWriter.CreateHeader(hdr)
		if err != nil {
			return fmt.Errorf("failed writing zip header: %s", err)
		}
		if fi.IsDir() {
			return nil
		}

		//nolint:gosec
		src, err := os.Open(fpath)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", fpath, err)
		}
		if _, err := io.Copy(w, src); err != nil {
			if err := src.Close(); err != nil {
				log.Println(err)
			}
			return fmt.Errorf("failed writing zip entry: %w", err)
		}
		if err := src.Close(); err != nil {
			log.Println(err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to traverse directory %q: %w", srcDir, err)
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed writing %q: %w", fpath, err)
	}
	if err := tgt.Close(); err != nil {
		return fmt.Errorf("failed writing %q: %w", fpath, err)
	}

	log.Printf("Successfully created %q", fpath)
	return nil
}

//nolint
func createTarball(srcDir, version, variantStr, sfx, grafanaDir string) error {
	fpath := filepath.Join(grafanaDir, "dist", fmt.Sprintf("grafana%s-%s.%s.tar.gz", sfx, version, variantStr))
	//nolint:gosec
	tgt, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", fpath, err)
	}
	defer func() {
		if err := tgt.Close(); err != nil {
			log.Println(err)
		}
	}()

	//nolint:gosec
	if err := os.Chmod(fpath, 0664); err != nil {
		return fmt.Errorf("failed to set permissions on %q: %w", fpath, err)
	}
	gzWriter, err := gzip.NewWriterLevel(tgt, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %w", err)
	}
	defer func() {
		if err := gzWriter.Close(); err != nil {
			log.Println(err)
		}
	}()
	tarWriter := tar.NewWriter(gzWriter)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, fname := range []string{"LICENSE", "README.md", "NOTICE.md"} {
		fpath := filepath.Join(grafanaDir, fname)
		fi, err := os.Lstat(fpath)
		if err != nil {
			return fmt.Errorf("couldn't stat %q: %w", fpath, err)
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return fmt.Errorf("failed getting tar header: %w", err)
		}
		hdr.Name = fmt.Sprintf("grafana-%s/%s", version, fname)
		if err := tarWriter.WriteHeader(hdr); err != nil {
			return fmt.Errorf("failed writing tar header: %w", err)
		}
		//nolint:gosec
		src, err := os.Open(fpath)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", fname, err)
		}
		if _, err := io.Copy(tarWriter, src); err != nil {
			if err := src.Close(); err != nil {
				log.Println(err)
			}
			return fmt.Errorf("failed writing tar entry: %w", err)
		}
		if err := src.Close(); err != nil {
			log.Println(err)
		}
	}
	if err := filepath.Walk(srcDir, func(fpath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fpath == srcDir {
			return nil
		}

		linkTgt := ""
		if fi.Mode()&os.ModeSymlink != 0 {
			log.Printf("reading link '%s'", fpath)
			linkTgt, err = os.Readlink(fpath)
			if err != nil {
				return err
			}
			linkTgt = fmt.Sprintf("grafana-%s/%s", version, linkTgt)
		}

		hdr, err := tar.FileInfoHeader(fi, linkTgt)
		if err != nil {
			return fmt.Errorf("failed getting tar header: %w", err)
		}
		hdr.Name = fmt.Sprintf("grafana-%s/%s", version, strings.TrimPrefix(fpath, fmt.Sprintf("%s/", srcDir)))
		if err := tarWriter.WriteHeader(hdr); err != nil {
			return fmt.Errorf("failed writing tar header: %w", err)
		}
		if fi.IsDir() {
			return nil
		}

		//nolint:gosec
		src, err := os.Open(fpath)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", fpath, err)
		}
		if _, err := io.Copy(tarWriter, src); err != nil {
			if err := src.Close(); err != nil {
				log.Println(err)
			}
			return fmt.Errorf("failed writing tar entry: %w", err)
		}
		if err := src.Close(); err != nil {
			log.Println(err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to traverse directory %q: %w", srcDir, err)
	}

	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("failed writing %q: %w", fpath, err)
	}
	if err := gzWriter.Close(); err != nil {
		return fmt.Errorf("failed writing %q: %w", fpath, err)
	}
	if err := tgt.Close(); err != nil {
		return fmt.Errorf("failed writing %q: %w", fpath, err)
	}

	st, err := os.Stat(fpath)
	if err != nil {
		return err
	}
	perms := st.Mode() & os.ModePerm
	log.Printf("Successfully created package %q (permissions: %o)", fpath, perms)

	return nil
}
// This file is autogenerated. DO NOT EDIT.
//
// Generated by pkg/framework/coremodel/gen.go
//
// Derived from the Thema lineage declared in pkg/coremodel/dashboard/coremodel.cue
//
// Run `make gen-cue` from repository root to regenerate.

package dashboard

import (
	"embed"
	"path/filepath"

	"github.com/grafana/grafana/pkg/cuectx"
	"github.com/grafana/grafana/pkg/framework/coremodel"
	"github.com/grafana/thema"
)

// Defines values for GraphTooltip.
const (
	GraphTooltipN0 GraphTooltip = 0

	GraphTooltipN1 GraphTooltip = 1

	GraphTooltipN2 GraphTooltip = 2
)

// Defines values for Style.
const (
	StyleDark Style = "dark"

	StyleLight Style = "light"
)

// Defines values for Timezone.
const (
	TimezoneBrowser Timezone = "browser"

	TimezoneEmpty Timezone = ""

	TimezoneUtc Timezone = "utc"
)

// Defines values for DashboardCursorSync.
const (
	DashboardCursorSyncN0 DashboardCursorSync = 0

	DashboardCursorSyncN1 DashboardCursorSync = 1

	DashboardCursorSyncN2 DashboardCursorSync = 2
)

// Defines values for DashboardLinkType.
const (
	DashboardLinkTypeDashboards DashboardLinkType = "dashboards"

	DashboardLinkTypeLink DashboardLinkType = "link"
)

// Defines values for FieldColorModeId.
const (
	FieldColorModeIdContinuousGrYlRd FieldColorModeId = "continuous-GrYlRd"

	FieldColorModeIdFixed FieldColorModeId = "fixed"

	FieldColorModeIdPaletteClassic FieldColorModeId = "palette-classic"

	FieldColorModeIdPaletteSaturated FieldColorModeId = "palette-saturated"

	FieldColorModeIdThresholds FieldColorModeId = "thresholds"
)

// Defines values for FieldColorSeriesByMode.
const (
	FieldColorSeriesByModeLast FieldColorSeriesByMode = "last"

	FieldColorSeriesByModeMax FieldColorSeriesByMode = "max"

	FieldColorSeriesByModeMin FieldColorSeriesByMode = "min"
)

// Defines values for GraphPanelType.
const (
	GraphPanelTypeGraph GraphPanelType = "graph"
)

// Defines values for HeatmapPanelType.
const (
	HeatmapPanelTypeHeatmap HeatmapPanelType = "heatmap"
)

// Defines values for PanelRepeatDirection.
const (
	PanelRepeatDirectionH PanelRepeatDirection = "h"

	PanelRepeatDirectionV PanelRepeatDirection = "v"
)

// Defines values for RowPanelType.
const (
	RowPanelTypeRow RowPanelType = "row"
)

// Defines values for ThresholdsConfigMode.
const (
	ThresholdsConfigModeAbsolute ThresholdsConfigMode = "absolute"

	ThresholdsConfigModePercentage ThresholdsConfigMode = "percentage"
)

// Defines values for ThresholdsMode.
const (
	ThresholdsModeAbsolute ThresholdsMode = "absolute"

	ThresholdsModePercentage ThresholdsMode = "percentage"
)

// Defines values for VariableModelType.
const (
	VariableModelTypeAdhoc VariableModelType = "adhoc"

	VariableModelTypeConstant VariableModelType = "constant"

	VariableModelTypeCustom VariableModelType = "custom"

	VariableModelTypeDatasource VariableModelType = "datasource"

	VariableModelTypeInterval VariableModelType = "interval"

	VariableModelTypeQuery VariableModelType = "query"

	VariableModelTypeSystem VariableModelType = "system"

	VariableModelTypeTextbox VariableModelType = "textbox"
)

// Defines values for VariableType.
const (
	VariableTypeAdhoc VariableType = "adhoc"

	VariableTypeConstant VariableType = "constant"

	VariableTypeCustom VariableType = "custom"

	VariableTypeDatasource VariableType = "datasource"

	VariableTypeInterval VariableType = "interval"

	VariableTypeQuery VariableType = "query"

	VariableTypeSystem VariableType = "system"

	VariableTypeTextbox VariableType = "textbox"
)

// Model is the Go representation of a dashboard.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Model struct {
	Annotations *struct {
		// TODO docs
		List []AnnotationQuery `json:"list"`
	} `json:"annotations,omitempty"`

	// Description of dashboard.
	Description *string `json:"description,omitempty"`

	// Whether a dashboard is editable or not.
	Editable bool `json:"editable"`

	// TODO docs
	FiscalYearStartMonth *int         `json:"fiscalYearStartMonth,omitempty"`
	GnetId               *string      `json:"gnetId,omitempty"`
	GraphTooltip         GraphTooltip `json:"graphTooltip"`

	// Unique numeric identifier for the dashboard.
	// TODO must isolate or remove identifiers local to a Grafana instance...?
	Id *int64 `json:"id,omitempty"`

	// TODO docs
	Links *[]DashboardLink `json:"links,omitempty"`

	// TODO docs
	LiveNow *bool          `json:"liveNow,omitempty"`
	Panels  *[]interface{} `json:"panels,omitempty"`

	// TODO docs
	Refresh *interface{} `json:"refresh,omitempty"`

	// Version of the JSON schema, incremented each time a Grafana update brings
	// changes to said schema.
	// TODO this is the existing schema numbering system. It will be replaced by Thema's themaVersion
	SchemaVersion int `json:"schemaVersion"`

	// Theme of dashboard.
	Style Style `json:"style"`

	// Tags associated with dashboard.
	Tags       *[]string `json:"tags,omitempty"`
	Templating *struct {
		// TODO docs
		List []VariableModel `json:"list"`
	} `json:"templating,omitempty"`

	// Time range for dashboard, e.g. last 6 hours, last 7 days, etc
	Time *struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"time,omitempty"`

	// TODO docs
	// TODO this appears to be spread all over in the frontend. Concepts will likely need tidying in tandem with schema changes
	Timepicker *struct {
		// Whether timepicker is collapsed or not.
		Collapse bool `json:"collapse"`

		// Whether timepicker is enabled or not.
		Enable bool `json:"enable"`

		// Whether timepicker is visible or not.
		Hidden bool `json:"hidden"`

		// Selectable intervals for auto-refresh.
		RefreshIntervals []string `json:"refresh_intervals"`

		// TODO docs
		TimeOptions []string `json:"time_options"`
	} `json:"timepicker,omitempty"`

	// Timezone of dashboard,
	Timezone *Timezone `json:"timezone,omitempty"`

	// Title of dashboard.
	Title *string `json:"title,omitempty"`

	// Unique dashboard identifier that can be generated by anyone. string (8-40)
	Uid *string `json:"uid,omitempty"`

	// Version of the dashboard, incremented each time the dashboard is updated.
	Version *int `json:"version,omitempty"`

	// TODO docs
	WeekStart *string `json:"weekStart,omitempty"`
}

// GraphTooltip is the Go representation of a Model.GraphTooltip.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type GraphTooltip int

// Theme of dashboard.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Style string

// Timezone of dashboard,
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Timezone string

// TODO docs
// FROM: AnnotationQuery in grafana-data/src/types/annotations.ts
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type AnnotationQuery struct {
	BuiltIn int `json:"builtIn"`

	// Datasource to use for annotation.
	Datasource struct {
		Type *string `json:"type,omitempty"`
		Uid  *string `json:"uid,omitempty"`
	} `json:"datasource"`

	// Whether annotation is enabled.
	Enable bool `json:"enable"`

	// Whether to hide annotation.
	Hide *bool `json:"hide,omitempty"`

	// Annotation icon color.
	IconColor *string `json:"iconColor,omitempty"`

	// Name of annotation.
	Name *string `json:"name,omitempty"`

	// Query for annotation data.
	RawQuery *string `json:"rawQuery,omitempty"`
	ShowIn   int     `json:"showIn"`

	// Schema for panel targets is specified by datasource
	// plugins. We use a placeholder definition, which the Go
	// schema loader either left open/as-is with the Base
	// variant of the Model and Panel families, or filled
	// with types derived from plugins in the Instance variant.
	// When working directly from CUE, importers can extend this
	// type directly to achieve the same effect.
	Target *Target `json:"target,omitempty"`
	Type   string  `json:"type"`
}

// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type DashboardCursorSync int

// FROM public/app/features/dashboard/state/Models.ts - ish
// TODO docs
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type DashboardLink struct {
	AsDropdown  bool              `json:"asDropdown"`
	Icon        *string           `json:"icon,omitempty"`
	IncludeVars bool              `json:"includeVars"`
	KeepTime    bool              `json:"keepTime"`
	Tags        []string          `json:"tags"`
	TargetBlank bool              `json:"targetBlank"`
	Title       string            `json:"title"`
	Tooltip     *string           `json:"tooltip,omitempty"`
	Type        DashboardLinkType `json:"type"`
	Url         *string           `json:"url,omitempty"`
}

// DashboardLinkType is the Go representation of a DashboardLink.Type.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type DashboardLinkType string

// TODO docs
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type FieldColor struct {
	// Stores the fixed color value if mode is fixed
	FixedColor *string `json:"fixedColor,omitempty"`

	// The main color scheme mode
	Mode interface{} `json:"mode"`

	// TODO docs
	SeriesBy *FieldColorSeriesByMode `json:"seriesBy,omitempty"`
}

// TODO docs
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type FieldColorModeId string

// TODO docs
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type FieldColorSeriesByMode string

// GraphPanel is the Go representation of a dashboard.GraphPanel.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type GraphPanel struct {
	// Support for legacy graph and heatmap panels.
	Type GraphPanelType `json:"type"`
}

// Support for legacy graph and heatmap panels.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type GraphPanelType string

// GridPos is the Go representation of a dashboard.GridPos.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type GridPos struct {
	// Panel
	H int `json:"h"`

	// true if fixed
	Static *bool `json:"static,omitempty"`

	// Panel
	W int `json:"w"`

	// Panel x
	X int `json:"x"`

	// Panel y
	Y int `json:"y"`
}

// HeatmapPanel is the Go representation of a dashboard.HeatmapPanel.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type HeatmapPanel struct {
	Type HeatmapPanelType `json:"type"`
}

// HeatmapPanelType is the Go representation of a HeatmapPanel.Type.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type HeatmapPanelType string

// Model panels. Panels are canonically defined inline
// because they share a version timeline with the dashboard
// schema; they do not evolve independently.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Panel struct {
	// The datasource used in all targets.
	Datasource *struct {
		Type *string `json:"type,omitempty"`
		Uid  *string `json:"uid,omitempty"`
	} `json:"datasource,omitempty"`

	// Description.
	Description *string `json:"description,omitempty"`
	FieldConfig struct {
		Defaults struct {
			// TODO docs
			Color *FieldColor `json:"color,omitempty"`

			// custom is specified by the PanelFieldConfig field
			// in panel plugin schemas.
			Custom *map[string]interface{} `json:"custom,omitempty"`

			// Significant digits (for display)
			Decimals *float32 `json:"decimals,omitempty"`

			// Human readable field metadata
			Description *string `json:"description,omitempty"`

			// The display value for this field.  This supports template variables blank is auto
			DisplayName *string `json:"displayName,omitempty"`

			// This can be used by data sources that return and explicit naming structure for values and labels
			// When this property is configured, this value is used rather than the default naming strategy.
			DisplayNameFromDS *string `json:"displayNameFromDS,omitempty"`

			// True if data source field supports ad-hoc filters
			Filterable *bool `json:"filterable,omitempty"`

			// // The behavior when clicking on a result
			Links *[]interface{} `json:"links,omitempty"`

			// Convert input values into a display string
			//
			// TODO this one corresponds to a complex type with
			// generics on the typescript side. Ouch. Will
			// either need special care, or we'll just need to
			// accept a very loosely specified schema. It's very
			// unlikely we'll be able to translate cue to
			// typescript generics in the general case, though
			// this particular one *may* be able to work.
			Mappings *[]map[string]interface{} `json:"mappings,omitempty"`
			Max      *float32                  `json:"max,omitempty"`
			Min      *float32                  `json:"min,omitempty"`

			// Alternative to empty string
			NoValue *string `json:"noValue,omitempty"`

			// An explict path to the field in the datasource.  When the frame meta includes a path,
			// This will default to `${frame.meta.path}/${field.name}
			//
			// When defined, this value can be used as an identifier within the datasource scope, and
			// may be used to update the results
			Path       *string           `json:"path,omitempty"`
			Thresholds *ThresholdsConfig `json:"thresholds,omitempty"`

			// Numeric Options
			Unit *string `json:"unit,omitempty"`

			// True if data source can write a value to the path.  Auth/authz are supported separately
			Writeable *bool `json:"writeable,omitempty"`
		} `json:"defaults"`
		Overrides []struct {
			Matcher struct {
				Id      string       `json:"id"`
				Options *interface{} `json:"options,omitempty"`
			} `json:"matcher"`
			Properties []struct {
				Id    string       `json:"id"`
				Value *interface{} `json:"value,omitempty"`
			} `json:"properties"`
		} `json:"overrides"`
	} `json:"fieldConfig"`
	GridPos *GridPos `json:"gridPos,omitempty"`

	// TODO docs
	Id *int `json:"id,omitempty"`

	// TODO docs
	// TODO tighter constraint
	Interval *string `json:"interval,omitempty"`

	// Panel links.
	// TODO fill this out - seems there are a couple variants?
	Links *[]DashboardLink `json:"links,omitempty"`

	// TODO docs
	MaxDataPoints *float32 `json:"maxDataPoints,omitempty"`

	// options is specified by the PanelOptions field in panel
	// plugin schemas.
	Options map[string]interface{} `json:"options"`

	// FIXME this almost certainly has to be changed in favor of scuemata versions
	PluginVersion *string `json:"pluginVersion,omitempty"`

	// Name of template variable to repeat for.
	Repeat *string `json:"repeat,omitempty"`

	// Direction to repeat in if 'repeat' is set.
	// "h" for horizontal, "v" for vertical.
	RepeatDirection PanelRepeatDirection `json:"repeatDirection"`

	// TODO docs
	Tags *[]string `json:"tags,omitempty"`

	// TODO docs
	Targets *[]Target `json:"targets,omitempty"`

	// TODO docs - seems to be an old field from old dashboard alerts?
	Thresholds *[]interface{} `json:"thresholds,omitempty"`

	// TODO docs
	// TODO tighter constraint
	TimeFrom *string `json:"timeFrom,omitempty"`

	// TODO docs
	TimeRegions *[]interface{} `json:"timeRegions,omitempty"`

	// TODO docs
	// TODO tighter constraint
	TimeShift *string `json:"timeShift,omitempty"`

	// Panel title.
	Title           *string `json:"title,omitempty"`
	Transformations []struct {
		Id      string                 `json:"id"`
		Options map[string]interface{} `json:"options"`
	} `json:"transformations"`

	// Whether to display the panel without a background.
	Transparent bool `json:"transparent"`

	// The panel plugin type id. May not be empty.
	Type string `json:"type"`
}

// Direction to repeat in if 'repeat' is set.
// "h" for horizontal, "v" for vertical.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type PanelRepeatDirection string

// Row panel
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type RowPanel struct {
	Collapsed bool `json:"collapsed"`

	// Name of default datasource.
	Datasource *struct {
		Type *string `json:"type,omitempty"`
		Uid  *string `json:"uid,omitempty"`
	} `json:"datasource,omitempty"`
	GridPos *GridPos      `json:"gridPos,omitempty"`
	Id      int           `json:"id"`
	Panels  []interface{} `json:"panels"`

	// Name of template variable to repeat for.
	Repeat *string      `json:"repeat,omitempty"`
	Title  *string      `json:"title,omitempty"`
	Type   RowPanelType `json:"type"`
}

// RowPanelType is the Go representation of a RowPanel.Type.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type RowPanelType string

// Schema for panel targets is specified by datasource
// plugins. We use a placeholder definition, which the Go
// schema loader either left open/as-is with the Base
// variant of the Model and Panel families, or filled
// with types derived from plugins in the Instance variant.
// When working directly from CUE, importers can extend this
// type directly to achieve the same effect.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Target map[string]interface{}

// TODO docs
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Threshold struct {
	// TODO docs
	Color string `json:"color"`

	// TODO docs
	// TODO are the values here enumerable into a disjunction?
	// Some seem to be listed in typescript comment
	State *string `json:"state,omitempty"`

	// TODO docs
	// FIXME the corresponding typescript field is required/non-optional, but nulls currently appear here when serializing -Infinity to JSON
	Value *float32 `json:"value,omitempty"`
}

// ThresholdsConfig is the Go representation of a dashboard.ThresholdsConfig.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type ThresholdsConfig struct {
	Mode ThresholdsConfigMode `json:"mode"`

	// Must be sorted by 'value', first value is always -Infinity
	Steps []struct {
		// TODO docs
		Color string `json:"color"`

		// TODO docs
		// TODO are the values here enumerable into a disjunction?
		// Some seem to be listed in typescript comment
		State *string `json:"state,omitempty"`

		// TODO docs
		// FIXME the corresponding typescript field is required/non-optional, but nulls currently appear here when serializing -Infinity to JSON
		Value *float32 `json:"value,omitempty"`
	} `json:"steps"`
}

// ThresholdsConfigMode is the Go representation of a ThresholdsConfig.Mode.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type ThresholdsConfigMode string

// ThresholdsMode is the Go representation of a dashboard.ThresholdsMode.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type ThresholdsMode string

// TODO docs
// FIXME this is extremely underspecfied; wasn't obvious which typescript types corresponded to it
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type Transformation struct {
	Id      string                 `json:"id"`
	Options map[string]interface{} `json:"options"`
}

// FROM: packages/grafana-data/src/types/templateVars.ts
// TODO docs
// TODO what about what's in public/app/features/types.ts?
// TODO there appear to be a lot of different kinds of [template] vars here? if so need a disjunction
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type VariableModel struct {
	Label *string           `json:"label,omitempty"`
	Name  string            `json:"name"`
	Type  VariableModelType `json:"type"`
}

// VariableModelType is the Go representation of a VariableModel.Type.
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type VariableModelType string

// FROM: packages/grafana-data/src/types/templateVars.ts
// TODO docs
// TODO this implies some wider pattern/discriminated union, probably?
//
// THIS TYPE IS INTENDED FOR INTERNAL USE BY THE GRAFANA BACKEND, AND IS SUBJECT TO BREAKING CHANGES.
// Equivalent Go types at stable import paths are provided in https://github.com/grafana/grok.
type VariableType string

//go:embed coremodel.cue
var cueFS embed.FS

// The current version of the coremodel schema, as declared in coremodel.cue.
// This version determines what schema version is returned from [Coremodel.CurrentSchema],
// and which schema version is used for code generation within the grafana/grafana repository.
//
// The code generator ensures that this is always the latest Thema schema version.
var currentVersion = thema.SV(0, 0)

// Lineage returns the Thema lineage representing a Grafana dashboard.
//
// The lineage is the canonical specification of the current dashboard schema,
// all prior schema versions, and the mappings that allow migration between
// schema versions.
func Lineage(lib thema.Library, opts ...thema.BindOption) (thema.Lineage, error) {
	return cuectx.LoadGrafanaInstancesWithThema(filepath.Join("pkg", "coremodel", "dashboard"), cueFS, lib, opts...)
}

var _ thema.LineageFactory = Lineage
var _ coremodel.Interface = &Coremodel{}

// Coremodel contains the foundational schema declaration for dashboards.
// It implements coremodel.Interface.
type Coremodel struct {
	lin thema.Lineage
}

// Lineage returns the canonical dashboard Lineage.
func (c *Coremodel) Lineage() thema.Lineage {
	return c.lin
}

// CurrentSchema returns the current (latest) dashboard Thema schema.
func (c *Coremodel) CurrentSchema() thema.Schema {
	return thema.SchemaP(c.lin, currentVersion)
}

// GoType returns a pointer to an empty Go struct that corresponds to
// the current Thema schema.
func (c *Coremodel) GoType() interface{} {
	return &Model{}
}

// New returns a new instance of the dashboard coremodel.
//
// Note that this function does not cache, and initially loading a Thema lineage
// can be expensive. As such, the Grafana backend should prefer to access this
// coremodel through a registry (pkg/framework/coremodel/registry), which does cache.
func New(lib thema.Library) (*Coremodel, error) {
	lin, err := Lineage(lib)
	if err != nil {
		return nil, err
	}

	return &Coremodel{
		lin: lin,
	}, nil
}
