package main

		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Templates []workflowTemplateStep `yaml:"templates"`
	} `yaml:"spec"`
}

type workflowTemplateStep struct {
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
		Env          []envVar `yaml:"env"`
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
}

type envVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type clusterWorkflow struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

// scriptForTemplate returns container.args[0] for the named Argo template.
func scriptForTemplate(t *testing.T, filename, templateName string) string {
	t.Helper()
	tmpl := workflowTemplateByName(t, filename, templateName)
	args := tmpl.Container.Args
	require.NotEmpty(t, args, "template %s/%s has no container.args", filename, templateName)
	return args[0]
}

func envForTemplate(t *testing.T, filename, templateName string) []envVar {
	t.Helper()
	tmpl := workflowTemplateByName(t, filename, templateName)
	return tmpl.Container.Env
}

func workflowTemplateByName(t *testing.T, filename, templateName string) workflowTemplateStep {
	t.Helper()
	wt := loadTemplate(t, filename)
	for _, tmpl := range wt.Spec.Templates {
		if tmpl.Name != templateName {
			continue
		}
		return tmpl
	}
	require.Failf(t, "template not found", "template %s not found in %s", templateName, filename)
	return workflowTemplateStep{}
}

func writeExec(t *testing.T, path, content string) {
echo "jq $*" >> "$CALLS"
echo "jq-stdin $input" >> "$CALLS"
case "$*" in
  *)
    case "$input" in
      *HTTP_PROXY*)
        printf '%s\n' "HTTP_PROXY=http://proxy"
        ;;
      *FOO*)
        printf '%s\n' "FOO=bar" "HELLO=world"
        ;;
    esac
    ;;
esac
exit 0
}

func runScript(t *testing.T, script string, stubs map[string]string, setup func(root string), replacements func(root string) []string) scriptRunResult {
	t.Helper()
	return runScriptWithEnv(t, script, nil, stubs, setup, replacements)
}

func runScriptWithEnv(t *testing.T, script string, templateEnv []envVar, stubs map[string]string, setup func(root string), replacements func(root string) []string) scriptRunResult {
	t.Helper()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping behavioral test")
		setup(root)
	}

	replacementPairs := []string(nil)
	if replacements != nil {
		replacementPairs = replacements(root)
		script = strings.NewReplacer(replacementPairs...).Replace(script)
	}

	env := append(os.Environ(),
		"PATH="+stubDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"CALLS="+callsFile,
	)
	if len(templateEnv) > 0 {
		replacer := strings.NewReplacer(replacementPairs...)
		for _, item := range templateEnv {
			env = append(env, item.Name+"="+replacer.Replace(item.Value))
		}
	}

	cmd := exec.Command("sh", "-c", script)
	cmd.Env = env
	out, err := cmd.CombinedOutput()

	res := scriptRunResult{output: string(out), root: root}

func TestContainerfileBuild_Behavior(t *testing.T) {
	script := scriptForTemplate(t, "containerfile-build.yaml", "build-image")
	env := envForTemplate(t, "containerfile-build.yaml", "build-image")
	res := runScriptWithEnv(t, script, env, map[string]string{
		"podman": buildPodmanStub,
		"jq":     buildJQStub,
	}, func(root string) {

func TestContainerfileBuild_MissingDockerfileFailsBeforeBuild(t *testing.T) {
	script := scriptForTemplate(t, "containerfile-build.yaml", "build-image")
	env := envForTemplate(t, "containerfile-build.yaml", "build-image")
	res := runScriptWithEnv(t, script, env, map[string]string{
		"podman": buildPodmanStub,
		"jq":     buildJQStub,
	}, func(root string) {
	} {
		t.Run(tc.file, func(t *testing.T) {
			script := scriptForTemplate(t, tc.file, "build-image")
			env := envForTemplate(t, tc.file, "build-image")
			res := runScriptWithEnv(t, script, env, map[string]string{
				"podman": buildPodmanStub,
				"pack":   packStub,
				"jq":     buildJQStub,

func TestBuildpackBuild_MissingAppPathFailsBeforePack(t *testing.T) {
	script := scriptForTemplate(t, "gcp-buildpacks-build.yaml", "build-image")
	env := envForTemplate(t, "gcp-buildpacks-build.yaml", "build-image")
	res := runScriptWithEnv(t, script, env, map[string]string{
		"podman": buildPodmanStub,
		"pack":   packStub,
		"jq":     buildJQStub,
	} {
		t.Run(tc.name, func(t *testing.T) {
			script := echoShim + scriptForTemplate(t, tc.file, "publish-image")
			env := envForTemplate(t, tc.file, "publish-image")
			res := runScriptWithEnv(t, script, env, map[string]string{"podman": publishPodmanStub}, func(root string) {
				require.NoError(t, os.MkdirAll(filepath.Join(root, "mnt-vol"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "storage"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "containers"), 0o755))
	} {
		t.Run(tc.name, func(t *testing.T) {
			script := echoShim + scriptForTemplate(t, tc.file, "publish-image")
			env := envForTemplate(t, tc.file, "publish-image")
			res := runScriptWithEnv(t, script, env, map[string]string{"podman": publishPodmanStub}, func(root string) {
				require.NoError(t, os.MkdirAll(filepath.Join(root, "mnt-vol"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "storage"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(root, "containers"), 0o755))

func TestPublishImage_MissingTarFailsBeforePush(t *testing.T) {
	script := scriptForTemplate(t, "publish-image.yaml", "publish-image")
	env := envForTemplate(t, "publish-image.yaml", "publish-image")
	res := runScriptWithEnv(t, script, env, map[string]string{"podman": publishPodmanStub}, func(root string) {
		require.NoError(t, os.MkdirAll(filepath.Join(root, "mnt-vol"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(root, "storage"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(root, "containers"), 0o755))
  %q`, contract, needle)
}

func requireEnvContains(t *testing.T, env []envVar, name string, valueNeedle string, contract string) {
	t.Helper()
	for _, item := range env {
		if item.Name == name && strings.Contains(item.Value, valueNeedle) {
			return
		}
	}
	t.Fatalf(`
contract:
  %s

expected env:
  %s contains %q

actual env:
  %v`, contract, name, valueNeedle, env)
}

// --- Cross-cutting invariants (scenario 21, 22) ---

func TestAllTemplates_ParseAndShape(t *testing.T) {
	for _, file := range buildTemplates {
		t.Run(file, func(t *testing.T) {
			s := scriptForTemplate(t, file, "build-image")
			env := envForTemplate(t, file, "build-image")
			// Output handoff to publish-image.
			requireContains(t, s, "/mnt/vol/app-image.tar")
			// Path validation guard.
			requireContains(t, s, "exit 1")
			// build-env JSON -> --env flags, with empty/[] skipped.
			requireContains(t, s, "BUILD_ENV_JSON", `!= "[]"`, "--env")
			requireEnvContains(t, env, "BUILD_ENV_JSON", "build-env",
				"build templates must receive build-env through container env, not raw shell interpolation")
			requireEnvContains(t, env, "IMAGE_NAME", "image-name",
				"build templates must receive image-name through container env")
			requireEnvContains(t, env, "IMAGE_TAG", "image-tag",
				"build templates must receive image-tag through container env")
			requireEnvContains(t, env, "GIT_REVISION", "git-revision",
				"build templates must receive git-revision through container env")
		})
	}
}

func TestContainerfileBuild_Specifics(t *testing.T) {
	s := scriptForTemplate(t, "containerfile-build.yaml", "build-image")
	env := envForTemplate(t, "containerfile-build.yaml", "build-image")
	requireContains(t, s,
		"podman build",
		"DOCKERFILE_PATH",
		"DOCKER_CONTEXT",
		"--build-arg", // containerfile additionally handles build-args
		"podman save -o /mnt/vol/app-image.tar",
	)
	requireEnvContains(t, env, "DOCKERFILE_PATH", "dockerfile-path",
		"containerfile build must receive dockerfile-path through container env")
	requireEnvContains(t, env, "DOCKER_CONTEXT", "docker-context",
		"containerfile build must receive docker-context through container env")
	requireEnvContains(t, env, "BUILD_ARGS_JSON", "build-args",
		"containerfile build must receive build-args through container env")
}

func TestBuildpackTemplates_Specifics(t *testing.T) {
	for _, file := range buildpackTemplates {
		t.Run(file, func(t *testing.T) {
			s := scriptForTemplate(t, file, "build-image")
			env := envForTemplate(t, file, "build-image")
			requireContains(t, s,
				"pack build",
				"--builder",
				"--run-image",
				"--pull-policy always",
				"--docker-host inherit",
				"APP_PATH",
			)
			requireEnvContains(t, env, "APP_PATH", "app-path",
				"buildpack templates must receive app-path through container env")
			// Supply-chain: builder/run images pinned by digest, not a tag.
			requireContains(t, s, "@sha256:")
			// Rootless podman service is started and waited on.
