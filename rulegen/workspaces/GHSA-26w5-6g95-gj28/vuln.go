package main


// GetWorkingDir returns the path to the workspace for this repo and pull.
func (w *FileWorkspace) GetWorkingDir(r models.Repo, p models.PullRequest, workspace string) (string, error) {
	repoDir := w.cloneDir(r, p, workspace)
	if _, err := os.Stat(repoDir); err != nil {
		return "", fmt.Errorf("checking if workspace exists: %w", err)
	}
// GetPullDir returns the dir where the workspaces for this pull are cloned.
// If the dir doesn't exist it will return an error.
func (w *FileWorkspace) GetPullDir(r models.Repo, p models.PullRequest) (string, error) {
	dir := w.repoPullDir(r, p)
	if _, err := os.Stat(dir); err != nil {
		return "", err
	}

// Delete deletes the workspace for this repo and pull.
func (w *FileWorkspace) Delete(logger logging.SimpleLogging, r models.Repo, p models.PullRequest) error {
	repoPullDir := w.repoPullDir(r, p)
	logger.Info("Deleting repo pull directory: %s", repoPullDir)
	return os.RemoveAll(repoPullDir)
}

// DeleteForWorkspace deletes the working dir for this workspace.
func (w *FileWorkspace) DeleteForWorkspace(logger logging.SimpleLogging, r models.Repo, p models.PullRequest, workspace string) error {
	workspaceDir := w.cloneDir(r, p, workspace)
	logger.Info("Deleting workspace directory: %s", workspaceDir)
	return os.RemoveAll(workspaceDir)
}

func (w *FileWorkspace) repoPullDir(r models.Repo, p models.PullRequest) string {
	return filepath.Join(w.DataDir, workingDirPrefix, r.FullName, strconv.Itoa(p.Num))
}

func (w *FileWorkspace) cloneDir(r models.Repo, p models.PullRequest, workspace string) string {
	return filepath.Join(w.repoPullDir(r, p), workspace)
}

// validateCloneDir computes the clone dir for the given repo/PR/workspace and
// pass through ATLANTIS_REPO_ALLOWLIST, but this explicit check makes the bound
// provable to static analyzers that taint-track user input into filesystem APIs.
func (w *FileWorkspace) validateCloneDir(r models.Repo, p models.PullRequest, workspace string) (string, error) {
	cloneDir := filepath.Clean(w.cloneDir(r, p, workspace))
	expectedPrefix := filepath.Clean(filepath.Join(w.DataDir, workingDirPrefix)) + string(filepath.Separator)
	if !strings.HasPrefix(cloneDir, expectedPrefix) {
		return "", fmt.Errorf("clone dir %q escapes managed data directory %q", cloneDir, expectedPrefix)
}

func (w *FileWorkspace) DeletePlan(logger logging.SimpleLogging, r models.Repo, p models.PullRequest, workspace string, projectPath string, projectName string) error {
	planPath := filepath.Join(w.cloneDir(r, p, workspace), projectPath, runtime.GetPlanFilename(workspace, projectName))
	logger.Info("Deleting plan: %s", planPath)
	return utils.RemoveIgnoreNonExistent(planPath)
}
// GitReadLock acquires a shared lock so that clone/reset/merge (write lock) cannot run
// while steps are using the working dir. Call the returned function when steps are done.
func (w *FileWorkspace) GitReadLock(r models.Repo, p models.PullRequest, workspace string) func() {
	return w.gitReadLock(w.cloneDir(r, p, workspace))
}

// gitReadLock acquires the same shared lock as GitReadLock but by workspace dir path.
	"github.com/runatlantis/atlantis/server/events/vcs"
	"github.com/runatlantis/atlantis/server/events/webhooks"
	"github.com/runatlantis/atlantis/server/logging"
)

const OperationComplete = true
		return nil, "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if _, err = os.Stat(absPath); os.IsNotExist(err) {

		// let's unlock here since something probably nuked our directory between the plan and policy check phase
	}

	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if _, err = os.Stat(projAbsPath); os.IsNotExist(err) {
		if unlockErr := lockAttempt.UnlockFn(); unlockErr != nil {
			ctx.Log.Err("error unlocking state after plan error: %v", unlockErr)
		return "", "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return "", "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
		return "", "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return "", "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
		return nil, "", cloneErr
	}
	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if _, err = os.Stat(projAbsPath); os.IsNotExist(err) {
		return nil, "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
		return nil, "", cloneErr
	}
	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if _, err = os.Stat(projAbsPath); os.IsNotExist(err) {
		return nil, "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
