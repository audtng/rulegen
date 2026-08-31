package main

	if err := rejectIsloSyncOptions(req); err != nil {
		return RunResult{}, err
	}
	started := b.now()
	client, err := newIsloClient(b.cfg, b.rt)
	if err != nil {
	fmt.Fprintf(b.rt.Stderr, "provider=islo lease=%s sandbox=%s\n", leaseID, name)
	syncDuration := time.Duration(0)
	syncPhases := []timingPhase{{Name: "sync", Skipped: true, Reason: "--no-sync"}}
	workspace := isloWorkspacePath(b.cfg)
	if !req.NoSync {
		var err error
		syncPhases, syncDuration, err = b.syncWorkspace(ctx, client, name, req)
}

func (b *isloBackend) createSandbox(ctx context.Context, client isloAPI, repo Repo, reclaim bool) (string, string, string, error) {
	name := newIsloSandboxName(repo)
	create := &gosdk.SandboxCreate{Name: stringValue(name)}
	if b.cfg.Islo.Image != "" {
		create.Image = stringValue(b.cfg.Islo.Image)
	}
	if b.cfg.Islo.Workdir != "" {
		create.Workdir = stringValue(b.cfg.Islo.Workdir)
	}
	if b.cfg.Islo.GatewayProfile != "" {
		create.GatewayProfile = stringValue(b.cfg.Islo.GatewayProfile)
	}
}

func TestIsloWorkspacePathDefaultsUnderWorkspace(t *testing.T) {
	if got := isloWorkspacePath(Config{}); got != "/workspace/crabbox" {
		t.Fatalf("workspace=%q", got)
	}
	if got := isloWorkspacePath(Config{Islo: IsloConfig{Workdir: "repo"}}); got != "/workspace/repo" {
		t.Fatalf("workspace=%q", got)
	}
	if got := isloWorkspacePath(Config{Islo: IsloConfig{Workdir: "/work/repo"}}); got != "/work/repo" {
		t.Fatalf("workspace=%q", got)
	}
}

	uploaded          bytes.Buffer
	uploadErr         error
	closeUploadReader bool
}

func (f *fakeIsloSyncClient) CreateSandbox(context.Context, *gosdk.SandboxCreate) (*gosdk.SandboxResponse, error) {
	return nil, nil
}

func (f *fakeIsloSyncClient) GetSandbox(context.Context, string) (*gosdk.SandboxResponse, error) {
		return nil, 0, err
	}
	preflightDuration := b.now().Sub(preflightStarted)
	workspace := isloWorkspacePath(b.cfg)
	prepareStarted := b.now()
	if err := b.prepareWorkspace(ctx, client, name, workspace); err != nil {
		return nil, 0, err
	return archive, nil
}

func isloWorkspacePath(cfg Config) string {
	workdir := strings.TrimSpace(cfg.Islo.Workdir)
	if workdir == "" {
		workdir = "crabbox"
	}
	if strings.HasPrefix(workdir, "/") {
		return path.Clean(workdir)
	}
	return path.Join("/workspace", workdir)
}
