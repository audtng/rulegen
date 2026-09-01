package main


	SourceDateEpoch time.Time
	WorkspaceDir    string
	WorkspaceDirFS  apkofs.FullFS
	WorkspaceIgnore string
	// Ordered directories where to find 'uses' pipelines.
	PipelineDirs          []string
	if b.Runner.Name() == container.QemuName {
		b.ExtraPackages = append(b.ExtraPackages, []string{
			"melange-microvm-init",
			"gnutar",
		}...)
	}

		defer os.RemoveAll(tmp)
		log.Infof("cache bucket copied to %s", tmp)

		fsys := apkofs.DirFS(tmp)

		// mkdir /var/cache/melange
		if err := os.MkdirAll(b.CacheDir, 0o755); err != nil {
		}

		log.Infof("populating workspace %s from %s", b.WorkspaceDir, b.SourceDir)
		if err := b.populateWorkspace(ctx, apkofs.DirFS(b.SourceDir)); err != nil {
			return fmt.Errorf("unable to populate workspace: %w", err)
		}
	}

	// Retrieve the post build workspace from the runner
	log.Infof("retrieving workspace from builder: %s", cfg.PodID)
	b.WorkspaceDirFS = apkofs.DirFS(b.WorkspaceDir)

	if err := b.retrieveWorkspace(ctx, b.WorkspaceDirFS); err != nil {
		return fmt.Errorf("retrieving workspace: %w", err)
	}
	log.Infof("retrieved and wrote post-build workspace to: %s", b.WorkspaceDir)
// filesystem in the directory `/var/lib/db/sbom`. The pkgName parameter should
// be set to the name of the origin package or subpackage.
func (b Build) writeSBOM(pkgName string, doc *spdx.Document) error {
	apkFSPath := filepath.Join(melangeOutputDirName, pkgName)
	sbomDirPath := filepath.Join(apkFSPath, "/var/lib/db/sbom")
	if err := b.WorkspaceDirFS.MkdirAll(sbomDirPath, os.FileMode(0o755)); err != nil {
		return fmt.Errorf("creating SBOM directory: %w", err)
	}

	pkgVersion := b.Configuration.Package.FullVersion()
	sbomPath := getPathForPackageSBOM(sbomDirPath, pkgName, pkgVersion)
	f, err := b.WorkspaceDirFS.Create(sbomPath)
	if err != nil {
		return fmt.Errorf("opening SBOM file for writing: %w", err)
	}
	"strings"
	"text/template"

	apkofs "chainguard.dev/apko/pkg/apk/fs"
	apko_types "chainguard.dev/apko/pkg/build/types"

	"github.com/klauspost/compress/gzip"
}

// TODO(kaniini): generate APKv3 packages
func (pc *PackageBuild) calculateInstalledSize(fsys apkofs.FullFS) error {
	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
	return nil
}

func (pc *PackageBuild) emitDataSection(ctx context.Context, fsys apkofs.FullFS, userinfofs apkofs.FullFS, remapUIDs map[int]int, remapGIDs map[int]int, w io.WriteSeeker) error {
	log := clog.FromContext(ctx)
	tarctx, err := tarball.NewContext(
		tarball.WithSourceDateEpoch(pc.Build.SourceDateEpoch),
	log.Info("generating package " + pc.Identity())

	// filesystem for the data package
	fsys, err := apkofs.Sub(pc.Build.WorkspaceDirFS, filepath.Join(melangeOutputDirName, pc.PackageName))
	if err != nil {
		return fmt.Errorf("failed to return filesystem for workspace subtree: %w", err)
	}

	// provide the tar writer etc/passwd and etc/group of guest filesystem
	userinfofs := apkofs.DirFS(pc.Build.GuestDir)

	hdl := &SCABuildInterface{
		PackageBuild: pc,
