package main

		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Templates []struct {
			Name   string `yaml:"name"`
			Inputs struct {
				Parameters []struct {
					Name    string `yaml:"name"`
					Default string `yaml:"default"`
				} `yaml:"parameters"`
			} `yaml:"inputs"`
			Container struct {
				Image        string   `yaml:"image"`
				Args         []string `yaml:"args"`
				VolumeMounts []struct {
					Name      string `yaml:"name"`
					MountPath string `yaml:"mountPath"`
					ReadOnly  bool   `yaml:"readOnly"`
				} `yaml:"volumeMounts"`
			} `yaml:"container"`
			Volumes []struct {
				Name   string `yaml:"name"`
				Secret *struct {
					SecretName string `yaml:"secretName"`
					Optional   *bool  `yaml:"optional"`
				} `yaml:"secret"`
			} `yaml:"volumes"`
		} `yaml:"templates"`
	} `yaml:"spec"`
}

type clusterWorkflow struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

// scriptForTemplate returns container.args[0] for the named Argo template.
func scriptForTemplate(t *testing.T, filename, templateName string) string {
	t.Helper()
	wt := loadTemplate(t, filename)
	for _, tmpl := range wt.Spec.Templates {
		if tmpl.Name != templateName {
			continue
		}
		args := tmpl.Container.Args
		require.NotEmpty(t, args, "template %s/%s has no container.args", filename, templateName)
		return args[0]
	}
	require.Failf(t, "template not found", "template %s not found in %s", templateName, filename)
	return ""
}

func writeExec(t *testing.T, path, content string) {
echo "jq $*" >> "$CALLS"
echo "jq-stdin $input" >> "$CALLS"
case "$*" in
  *"--env"*)
    printf '%s\n' "--env FOO=bar" "--env HELLO=world"
    ;;
  *"--build-arg"*)
    printf '%s\n' "--build-arg HTTP_PROXY=http://proxy"
    ;;
esac
exit 0
}

func runScript(t *testing.T, script string, stubs map[string]string, setup func(root string), replacements func(root string) []string) scriptRunResult {
	t.Helper()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping behavioral test")
		setup(root)
	}

	if replacements != nil {
		script = strings.NewReplacer(replacements(root)...).Replace(script)
	}
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = append(os.Environ(),
		"PATH="+stubDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"CALLS="+callsFile,
	)
	out, err := cmd.CombinedOutput()

	res := scriptRunResult{output: string(out), root: root}

func TestContainerfileBuild_Behavior(t *testing.T) {
	script := scriptForTemplate(t, "containerfile-build.yaml", "build-image")
	res := runScript(t, script, map[string]string{
		"podman": buildPodmanStub,
		"jq":     buildJQStub,
	}, func(root string) {

func TestContainerfileBuild_MissingDockerfileFailsBeforeBuild(t *testing.T) {
	script := scriptForTemplate(t, "containerfile-build.yaml", "build-image")
	res := runScript(t, script, map[string]string{
		"podman": buildPodmanStub,
		"jq":     buildJQStub,
	}, func(root string) {
	} {
		t.Run(tc.file, func(t *testing.T) {
			script := scriptForTemplate(t, tc.file, "build-image")
			res := runScript(t, script, map[string]string{
				"podman": buildPodmanStub,
				"pack":   packStub,
				"jq":     buildJQStub,

func TestBuildpackBuild_MissingAppPathFailsBeforePack(t *testing.T) {
	script := scriptForTemplate(t, "gcp-buildpacks-build.yaml", "build-image")
	res := runScript(t, script, map[string]string{
		"podman": buildPodmanStub,
		"pack":   packStub,
		"jq":     buildJQStub,
	} {
		t.Run(tc.name, func(t *testing.T) {
			script := echoShim + scriptForTemplate(t, tc.file, "publish-image")
			res := runScript(t, script, map[string]string{"podman": publishPodmanStub}, func(root string) {
				require.NoError(t, os.MkdirAll(filepath.Join(root, "mnt-vol"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "storage"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "containers"), 0o755))
	} {
		t.Run(tc.name, func(t *testing.T) {
			script := echoShim + scriptForTemplate(t, tc.file, "publish-image")
			res := runScript(t, script, map[string]string{"podman": publishPodmanStub}, func(root string) {
				require.NoError(t, os.MkdirAll(filepath.Join(root, "mnt-vol"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "storage"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "containers"), 0o755))

func TestPublishImage_MissingTarFailsBeforePush(t *testing.T) {
	script := scriptForTemplate(t, "publish-image.yaml", "publish-image")
	res := runScript(t, script, map[string]string{"podman": publishPodmanStub}, func(root string) {
		require.NoError(t, os.MkdirAll(filepath.Join(root, "mnt-vol"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(root, "storage"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(root, "containers"), 0o755))
  %q`, contract, needle)
}

// --- Cross-cutting invariants (scenario 21, 22) ---

func TestAllTemplates_ParseAndShape(t *testing.T) {
	for _, file := range buildTemplates {
		t.Run(file, func(t *testing.T) {
			s := scriptForTemplate(t, file, "build-image")
			// Output handoff to publish-image.
			requireContains(t, s, "/mnt/vol/app-image.tar")
			// Path validation guard.
			requireContains(t, s, "exit 1")
			// build-env JSON -> --env flags, with empty/[] skipped.
			requireContains(t, s, "build-env", `!= "[]"`, "--env")
		})
	}
}

func TestContainerfileBuild_Specifics(t *testing.T) {
	s := scriptForTemplate(t, "containerfile-build.yaml", "build-image")
	requireContains(t, s,
		"podman build",
		"dockerfile-path",
		"docker-context",
		"--build-arg", // containerfile additionally handles build-args
		"podman save -o /mnt/vol/app-image.tar",
	)
}

func TestBuildpackTemplates_Specifics(t *testing.T) {
	for _, file := range buildpackTemplates {
		t.Run(file, func(t *testing.T) {
			s := scriptForTemplate(t, file, "build-image")
			requireContains(t, s,
				"pack build",
				"--builder",
				"--run-image",
				"--pull-policy always",
				"--docker-host inherit",
				"app-path",
			)
			// Supply-chain: builder/run images pinned by digest, not a tag.
			requireContains(t, s, "@sha256:")
			// Rootless podman service is started and waited on.
