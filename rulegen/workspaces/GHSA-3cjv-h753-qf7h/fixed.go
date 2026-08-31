package main

	if err := rejectIsloSyncOptions(req); err != nil {
		return RunResult{}, err
	}
	workspace, err := isloWorkspacePath(b.cfg)
	if err != nil {
		return RunResult{}, err
	}
	started := b.now()
	client, err := newIsloClient(b.cfg, b.rt)
	if err != nil {
	fmt.Fprintf(b.rt.Stderr, "provider=islo lease=%s sandbox=%s\n", leaseID, name)
	syncDuration := time.Duration(0)
	syncPhases := []timingPhase{{Name: "sync", Skipped: true, Reason: "--no-sync"}}
	if !req.NoSync {
		var err error
		syncPhases, syncDuration, err = b.syncWorkspace(ctx, client, name, req)
}

func (b *isloBackend) createSandbox(ctx context.Context, client isloAPI, repo Repo, reclaim bool) (string, string, string, error) {
	workdir, err := isloRelativeWorkdir(b.cfg)
	if err != nil {
		return "", "", "", err
	}
	name := newIsloSandboxName(repo)
	create := &gosdk.SandboxCreate{Name: stringValue(name)}
	if b.cfg.Islo.Image != "" {
		create.Image = stringValue(b.cfg.Islo.Image)
	}
	create.Workdir = stringValue(workdir)
	if b.cfg.Islo.GatewayProfile != "" {
		create.GatewayProfile = stringValue(b.cfg.Islo.GatewayProfile)
	}
}

func TestIsloWorkspacePathDefaultsUnderWorkspace(t *testing.T) {
	if got, err := isloWorkspacePath(Config{}); err != nil || got != "/workspace/crabbox" {
		t.Fatalf("workspace=%q err=%v", got, err)
	}
	if got, err := isloWorkspacePath(Config{Islo: IsloConfig{Workdir: "repo"}}); err != nil || got != "/workspace/repo" {
		t.Fatalf("workspace=%q err=%v", got, err)
	}
	if got, err := isloWorkspacePath(Config{Islo: IsloConfig{Workdir: "team/repo"}}); err != nil || got != "/workspace/team/repo" {
		t.Fatalf("workspace=%q err=%v", got, err)
	}
}

func TestIsloWorkspacePathRejectsEscapes(t *testing.T) {
	for _, workdir := range []string{"/work/repo", "/etc", "../etc", "repo/../../../etc", ".", "./.."} {
		t.Run(workdir, func(t *testing.T) {
			if got, err := isloWorkspacePath(Config{Islo: IsloConfig{Workdir: workdir}}); err == nil {
				t.Fatalf("workspace=%q, want error for workdir %q", got, workdir)
			}
		})
	}
}

func TestIsloRunRejectsUnsafeWorkdirBeforeProviderClient(t *testing.T) {
	backend := &isloBackend{
		cfg: Config{Islo: IsloConfig{Workdir: "../etc"}},
		rt:  Runtime{Stderr: io.Discard},
	}
	_, err := backend.Run(context.Background(), RunRequest{NoSync: true})
	if err == nil || !strings.Contains(err.Error(), "escapes /workspace") {
		t.Fatalf("Run err=%v, want workdir containment error", err)
	}
}

func TestIsloCreateSandboxRejectsUnsafeWorkdirBeforeAPI(t *testing.T) {
	client := &fakeIsloSyncClient{}
	backend := &isloBackend{
		cfg: Config{Islo: IsloConfig{Workdir: "../etc"}},
		rt:  Runtime{Stderr: io.Discard},
	}
	_, _, _, err := backend.createSandbox(context.Background(), client, Repo{Root: t.TempDir(), Name: "repo"}, false)
	if err == nil || !strings.Contains(err.Error(), "escapes /workspace") {
		t.Fatalf("createSandbox err=%v, want workdir containment error", err)
	}
	if client.createRequest != nil {
		t.Fatalf("CreateSandbox was called with %#v", client.createRequest)
	}
}

func TestIsloCreateSandboxPassesRelativeWorkdirToProvider(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	client := &fakeIsloSyncClient{createName: "crabbox-repo-abcdef"}
	backend := &isloBackend{
		cfg: Config{Islo: IsloConfig{Workdir: "team/repo"}},
		rt:  Runtime{Stderr: io.Discard},
	}
	_, _, _, err := backend.createSandbox(context.Background(), client, Repo{Root: t.TempDir(), Name: "repo"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if client.createRequest == nil || client.createRequest.Workdir == nil || *client.createRequest.Workdir != "team/repo" {
		t.Fatalf("create workdir=%v", client.createRequest)
	}
}

	uploaded          bytes.Buffer
	uploadErr         error
	closeUploadReader bool
	createRequest     *gosdk.SandboxCreate
	createName        string
}

func (f *fakeIsloSyncClient) CreateSandbox(_ context.Context, req *gosdk.SandboxCreate) (*gosdk.SandboxResponse, error) {
	f.createRequest = req
	name := f.createName
	if name == "" {
		name = "crabbox-test-abcdef"
	}
	return &gosdk.SandboxResponse{Name: name}, nil
}

func (f *fakeIsloSyncClient) GetSandbox(context.Context, string) (*gosdk.SandboxResponse, error) {
		return nil, 0, err
	}
	preflightDuration := b.now().Sub(preflightStarted)
	workspace, err := isloWorkspacePath(b.cfg)
	if err != nil {
		return nil, 0, err
	}
	prepareStarted := b.now()
	if err := b.prepareWorkspace(ctx, client, name, workspace); err != nil {
		return nil, 0, err
	return archive, nil
}

func isloWorkspacePath(cfg Config) (string, error) {
	workdir, err := isloRelativeWorkdir(cfg)
	if err != nil {
		return "", err
	}
	return path.Join("/workspace", workdir), nil
}

func isloRelativeWorkdir(cfg Config) (string, error) {
	workdir := strings.TrimSpace(cfg.Islo.Workdir)
	if workdir == "" {
		workdir = "crabbox"
	}
	if strings.HasPrefix(workdir, "/") {
		return "", exit(2, "islo workdir %q must be relative under /workspace", workdir)
	}
	workdir = path.Clean(workdir)
	if workdir == "." || workdir == ".." || strings.HasPrefix(workdir, "../") {
		return "", exit(2, "islo workdir %q escapes /workspace", workdir)
	}
	return workdir, nil
}
