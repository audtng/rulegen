package main

		return nil
	}

	validWorkspace := func(value any) error {
		strPtr := value.(*string)
		if strPtr == nil || *strPtr == "" {
			return nil
		}
		ws := *strPtr
		if strings.Contains(ws, "..") || strings.ContainsAny(ws, "/\\") {
			return errors.New("cannot contain '..', '/', or '\\'")
		}
		return nil
	}

	// Validate that name doesn't contain glob patterns - glob expansion only works for 'dir'
	if p.Name != nil && ContainsGlobPattern(*p.Name) {
		return errors.New("name: cannot contain glob pattern characters ('*', '?', '['); glob expansion is only supported in the 'dir' field")

	return validation.ValidateStruct(&p,
		validation.Field(&p.Dir, validation.Required, validation.By(validDir)),
		validation.Field(&p.Workspace, validation.By(validWorkspace)),
		validation.Field(&p.PlanRequirements, validation.By(validPlanReq)),
		validation.Field(&p.ApplyRequirements, validation.By(validApplyReq)),
		validation.Field(&p.ImportRequirements, validation.By(validImportReq)),
			},
			expErr: "",
		},
		{
			description: "workspace with ..",
			input: raw.Project{
				Dir:       String("."),
				Workspace: String("../evil"),
			},
			expErr: "workspace: cannot contain '..', '/', or '\\'.",
		},
		{
			description: "workspace beginning with /",
			input: raw.Project{
				Dir:       String("."),
				Workspace: String("/etc"),
			},
			expErr: "workspace: cannot contain '..', '/', or '\\'.",
		},
		{
			description: "workspace with embedded /",
			input: raw.Project{
				Dir:       String("."),
				Workspace: String("sub/dir"),
			},
			expErr: "workspace: cannot contain '..', '/', or '\\'.",
		},
		{
			description: "workspace with backslash",
			input: raw.Project{
				Dir:       String("."),
				Workspace: String("sub\\dir"),
			},
			expErr: "workspace: cannot contain '..', '/', or '\\'.",
		},
		{
			description: "valid workspace",
			input: raw.Project{
				Dir:       String("."),
				Workspace: String("my-workspace"),
			},
			expErr: "",
		},
	}
	validation.ErrorTag = "yaml"
	for _, c := range cases {
	if strings.Contains(repo, "/") {
		return Repo{}, fmt.Errorf("invalid repo format %q, repo %q should not contain any /'s", repoFullName, owner)
	}
	if strings.Contains(owner, "..") || strings.Contains(repo, "..") {
		return Repo{}, fmt.Errorf("invalid repo format %q, owner or repo cannot contain '..'", repoFullName)
	}

	return Repo{
		FullName:          repoFullName,
			"/b",
			`invalid repo format "/b", owner "" or repo "b" was empty`,
		},
		{
			"owner../repo",
			`invalid repo format "owner../repo", owner or repo cannot contain '..'`,
		},
		{
			"owner/..repo",
			`invalid repo format "owner/..repo", owner or repo cannot contain '..'`,
		},
		// Trailing ".." in repo name
		{
			"owner/repo..",
			`invalid repo format "owner/repo..", owner or repo cannot contain '..'`,
		},
	}
	for _, c := range cases {
		t.Run(c.repoFullName, func(t *testing.T) {
			ErrEquals(t, c.expErr, err)
		})
	}

	// GitLab allows subgroups (slashes in owner), so the ".." check fires for
	// a path like "group/subgroup/../../../etc/repo" where the owner contains "..".
	t.Run("gitlab subgroup owner with ..", func(t *testing.T) {
		repoFullName := "group/subgroup/../../etc/repo"
		cloneURL := fmt.Sprintf("https://gitlab.com/%s.git", repoFullName)
		_, err := models.NewRepo(models.Gitlab, repoFullName, cloneURL, "u", "p", "")
		ErrEquals(t, `invalid repo format "group/subgroup/../../etc/repo", owner or repo cannot contain '..'`, err)
	})
}

// If the clone url doesn't end with .git, and VCS is not Azure DevOps, it is appended
	"github.com/runatlantis/atlantis/server/events/vcs"
	"github.com/runatlantis/atlantis/server/events/webhooks"
	"github.com/runatlantis/atlantis/server/logging"
	"github.com/runatlantis/atlantis/server/utils"
)

const OperationComplete = true
		return nil, "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if err := utils.EnsureSubPath(repoDir, absPath); err != nil {

		// let's unlock here since something probably nuked our directory between the plan and policy check phase
		if unlockErr := lockAttempt.UnlockFn(); unlockErr != nil {
			ctx.Log.Err("error unlocking state after plan error: %v", unlockErr)
		}

		return nil, "", fmt.Errorf("project path traversal detected: %w", err)
	}
	if _, err = os.Stat(absPath); os.IsNotExist(err) {

		// let's unlock here since something probably nuked our directory between the plan and policy check phase
	}

	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if err := utils.EnsureSubPath(repoDir, projAbsPath); err != nil {
		if unlockErr := lockAttempt.UnlockFn(); unlockErr != nil {
			ctx.Log.Err("error unlocking state after plan error: %v", unlockErr)
		}
		return nil, "", fmt.Errorf("project path traversal detected: %w", err)
	}
	if _, err = os.Stat(projAbsPath); os.IsNotExist(err) {
		if unlockErr := lockAttempt.UnlockFn(); unlockErr != nil {
			ctx.Log.Err("error unlocking state after plan error: %v", unlockErr)
		return "", "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if err := utils.EnsureSubPath(repoDir, absPath); err != nil {
		return "", "", fmt.Errorf("project path traversal detected: %w", err)
	}
	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return "", "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
		return "", "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if err := utils.EnsureSubPath(repoDir, absPath); err != nil {
		return "", "", fmt.Errorf("project path traversal detected: %w", err)
	}
	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return "", "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
		return nil, "", cloneErr
	}
	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if err = utils.EnsureSubPath(repoDir, projAbsPath); err != nil {
		return nil, "", fmt.Errorf("project path traversal detected: %w", err)
	}
	if _, err = os.Stat(projAbsPath); os.IsNotExist(err) {
		return nil, "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
		return nil, "", cloneErr
	}
	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)
	if err = utils.EnsureSubPath(repoDir, projAbsPath); err != nil {
		return nil, "", fmt.Errorf("project path traversal detected: %w", err)
	}
	if _, err = os.Stat(projAbsPath); os.IsNotExist(err) {
		return nil, "", DirNotExistErr{RepoRelDir: ctx.RepoRelDir}
	}
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
		})
	}
}

func TestDefaultProjectCommandRunner_PathTraversal(t *testing.T) {
	const defaultTraversalPattern = "../../../../etc"

	cases := []struct {
		name              string
		traversalPatterns []string
		expectUnlock      bool
		runFn             func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error
		setupWorkingDir   func(mockWorkingDir *mocks.MockWorkingDir, repoDir string)
	}{
		{
			name: "Plan",
			traversalPatterns: []string{
				"../../../../etc",
				"../etc",
				"sub/../../etc",
			},
			expectUnlock: true,
			setupWorkingDir: func(mockWorkingDir *mocks.MockWorkingDir, repoDir string) {
				When(mockWorkingDir.Clone(Any[logging.SimpleLogging](), Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(repoDir, nil)
				When(mockWorkingDir.MergeAgain(Any[logging.SimpleLogging](), Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(false, nil)
			},
			runFn: func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error {
				return runner.Plan(ctx).Error
			},
		},
		{
			name:              "Apply",
			traversalPatterns: []string{defaultTraversalPattern},
			setupWorkingDir: func(mockWorkingDir *mocks.MockWorkingDir, repoDir string) {
				When(mockWorkingDir.GetWorkingDir(Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(repoDir, nil)
			},
			runFn: func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error {
				return runner.Apply(ctx).Error
			},
		},
		{
			name:              "PolicyCheck",
			traversalPatterns: []string{defaultTraversalPattern},
			expectUnlock:      true,
			setupWorkingDir: func(mockWorkingDir *mocks.MockWorkingDir, repoDir string) {
				When(mockWorkingDir.GetWorkingDir(Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(repoDir, nil)
			},
			runFn: func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error {
				return runner.PolicyCheck(ctx).Error
			},
		},
		{
			name:              "Version",
			traversalPatterns: []string{defaultTraversalPattern},
			setupWorkingDir: func(mockWorkingDir *mocks.MockWorkingDir, repoDir string) {
				When(mockWorkingDir.GetWorkingDir(Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(repoDir, nil)
			},
			runFn: func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error {
				return runner.Version(ctx).Error
			},
		},
		{
			name:              "Import",
			traversalPatterns: []string{defaultTraversalPattern},
			setupWorkingDir: func(mockWorkingDir *mocks.MockWorkingDir, repoDir string) {
				When(mockWorkingDir.Clone(Any[logging.SimpleLogging](), Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(repoDir, nil)
			},
			runFn: func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error {
				return runner.Import(ctx).Error
			},
		},
		{
			name:              "StateRm",
			traversalPatterns: []string{defaultTraversalPattern},
			setupWorkingDir: func(mockWorkingDir *mocks.MockWorkingDir, repoDir string) {
				When(mockWorkingDir.Clone(Any[logging.SimpleLogging](), Any[models.Repo](), Any[models.PullRequest](), Any[string]())).
					ThenReturn(repoDir, nil)
			},
			runFn: func(runner *events.DefaultProjectCommandRunner, ctx command.ProjectContext) error {
				return runner.StateRm(ctx).Error
			},
		},
	}

	for _, tc := range cases {
		for _, pattern := range tc.traversalPatterns {
			t.Run(tc.name+" rejects traversal pattern "+pattern, func(t *testing.T) {
				RegisterMockTestingT(t)
				mockWorkingDir := mocks.NewMockWorkingDir()
				mockLocker := mocks.NewMockProjectLocker()
				runner := &events.DefaultProjectCommandRunner{
					Locker:           mockLocker,
					LockURLGenerator: mockURLGenerator{},
					WorkingDir:       mockWorkingDir,
					WorkingDirLocker: events.NewDefaultWorkingDirLocker(),
				}
				repoDir := t.TempDir()
				tc.setupWorkingDir(mockWorkingDir, repoDir)
				unlockCalled := false
				When(mockLocker.TryLock(
					Any[logging.SimpleLogging](),
					Any[models.PullRequest](),
					Any[models.User](),
					Any[string](),
					Any[models.Project](),
					AnyBool(),
				)).ThenReturn(&events.TryLockResponse{
					LockAcquired: true,
					LockKey:      "lock-key",
					UnlockFn: func() error {
						unlockCalled = true
						return nil
					},
				}, nil)
				ctx := command.ProjectContext{
					Log:        logging.NewNoopLogger(t),
					Workspace:  "default",
					RepoRelDir: pattern,
					RePlanCmd:  "atlantis plan -d .",
				}
				err := tc.runFn(runner, ctx)
				Assert(t, err != nil, "expected error for RepoRelDir %q in runner %q", pattern, tc.name)
				Assert(t,
					strings.Contains(err.Error(), "project path traversal detected"),
					"expected traversal error for runner %q with RepoRelDir %q, got: %s", tc.name, pattern, err,
				)
				if tc.expectUnlock {
					Assert(t, unlockCalled, "expected runner %q with RepoRelDir %q to release project lock", tc.name, pattern)
				}
			})
		}
	}
}

// GetWorkingDir returns the path to the workspace for this repo and pull.
func (w *FileWorkspace) GetWorkingDir(r models.Repo, p models.PullRequest, workspace string) (string, error) {
	repoDir, err := w.cloneDir(r, p, workspace)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(repoDir); err != nil {
		return "", fmt.Errorf("checking if workspace exists: %w", err)
	}
// GetPullDir returns the dir where the workspaces for this pull are cloned.
// If the dir doesn't exist it will return an error.
func (w *FileWorkspace) GetPullDir(r models.Repo, p models.PullRequest) (string, error) {
	dir, err := w.repoPullDir(r, p)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); err != nil {
		return "", err
	}

// Delete deletes the workspace for this repo and pull.
func (w *FileWorkspace) Delete(logger logging.SimpleLogging, r models.Repo, p models.PullRequest) error {
	repoPullDir, err := w.repoPullDir(r, p)
	if err != nil {
		return err
	}
	logger.Info("Deleting repo pull directory: %s", repoPullDir)
	return os.RemoveAll(repoPullDir)
}

// DeleteForWorkspace deletes the working dir for this workspace.
func (w *FileWorkspace) DeleteForWorkspace(logger logging.SimpleLogging, r models.Repo, p models.PullRequest, workspace string) error {
	workspaceDir, err := w.cloneDir(r, p, workspace)
	if err != nil {
		return err
	}
	logger.Info("Deleting workspace directory: %s", workspaceDir)
	return os.RemoveAll(workspaceDir)
}

func (w *FileWorkspace) repoPullDir(r models.Repo, p models.PullRequest) (string, error) {
	dir := filepath.Join(w.DataDir, workingDirPrefix, r.FullName, strconv.Itoa(p.Num))
	if err := utils.EnsureSubPath(filepath.Join(w.DataDir, workingDirPrefix), dir); err != nil {
		return "", fmt.Errorf("repo path traversal detected: %w", err)
	}
	return dir, nil
}

func (w *FileWorkspace) cloneDir(r models.Repo, p models.PullRequest, workspace string) (string, error) {
	pullDir, err := w.repoPullDir(r, p)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(pullDir, workspace)
	if err := utils.EnsureSubPath(pullDir, dir); err != nil {
		return "", fmt.Errorf("workspace path traversal detected: %w", err)
	}
	return dir, nil
}

// validateCloneDir computes the clone dir for the given repo/PR/workspace and
// pass through ATLANTIS_REPO_ALLOWLIST, but this explicit check makes the bound
// provable to static analyzers that taint-track user input into filesystem APIs.
func (w *FileWorkspace) validateCloneDir(r models.Repo, p models.PullRequest, workspace string) (string, error) {
	cloneDir, err := w.cloneDir(r, p, workspace)
	if err != nil {
		return "", err
	}
	cloneDir = filepath.Clean(cloneDir)
	expectedPrefix := filepath.Clean(filepath.Join(w.DataDir, workingDirPrefix)) + string(filepath.Separator)
	if !strings.HasPrefix(cloneDir, expectedPrefix) {
		return "", fmt.Errorf("clone dir %q escapes managed data directory %q", cloneDir, expectedPrefix)
}

func (w *FileWorkspace) DeletePlan(logger logging.SimpleLogging, r models.Repo, p models.PullRequest, workspace string, projectPath string, projectName string) error {
	cloneDir, err := w.cloneDir(r, p, workspace)
	if err != nil {
		return err
	}
	planPath := filepath.Join(cloneDir, projectPath, runtime.GetPlanFilename(workspace, projectName))
	if err := utils.EnsureSubPath(cloneDir, planPath); err != nil {
		return fmt.Errorf("plan path traversal detected: %w", err)
	}
	logger.Info("Deleting plan: %s", planPath)
	return utils.RemoveIgnoreNonExistent(planPath)
}
// GitReadLock acquires a shared lock so that clone/reset/merge (write lock) cannot run
// while steps are using the working dir. Call the returned function when steps are done.
func (w *FileWorkspace) GitReadLock(r models.Repo, p models.PullRequest, workspace string) func() {
	cloneDir, err := w.cloneDir(r, p, workspace)
	if err != nil {
		return func() {}
	}
	return w.gitReadLock(cloneDir)
}

// gitReadLock acquires the same shared lock as GitReadLock but by workspace dir path.
	return repoDir
}

func TestFileWorkspace_PathTraversal(t *testing.T) {
	logger := logging.NewNoopLogger(t)
	dataDir := t.TempDir()

	wd := &events.FileWorkspace{
		DataDir:             dataDir,
		GpgNoSigningEnabled: true,
	}

	maliciousRepo := models.Repo{FullName: "../../../etc"}
	pull := models.PullRequest{
		Num:      1,
		BaseRepo: maliciousRepo,
	}

	t.Run("GetWorkingDir rejects traversal in repo name", func(t *testing.T) {
		_, err := wd.GetWorkingDir(maliciousRepo, pull, "default")
		Assert(t, err != nil, "expected error for path traversal in repo name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("GetPullDir rejects traversal in repo name", func(t *testing.T) {
		_, err := wd.GetPullDir(maliciousRepo, pull)
		Assert(t, err != nil, "expected error for path traversal in repo name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("Delete rejects traversal in repo name", func(t *testing.T) {
		err := wd.Delete(logger, maliciousRepo, pull)
		Assert(t, err != nil, "expected error for path traversal in repo name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("DeleteForWorkspace rejects traversal in repo name", func(t *testing.T) {
		err := wd.DeleteForWorkspace(logger, maliciousRepo, pull, "default")
		Assert(t, err != nil, "expected error for path traversal in repo name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("DeleteForWorkspace rejects traversal in workspace name", func(t *testing.T) {
		safeRepo := models.Repo{FullName: "owner/repo"}
		safePull := models.PullRequest{Num: 1, BaseRepo: safeRepo}
		err := wd.DeleteForWorkspace(logger, safeRepo, safePull, "../../etc")
		Assert(t, err != nil, "expected error for path traversal in workspace name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("Clone rejects traversal in repo name", func(t *testing.T) {
		_, err := wd.Clone(logger, maliciousRepo, pull, "default")
		Assert(t, err != nil, "expected error for path traversal in repo name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("DeletePlan rejects traversal in project path", func(t *testing.T) {
		safeRepo := models.Repo{FullName: "owner/repo"}
		safePull := models.PullRequest{Num: 1, BaseRepo: safeRepo}
		cloneDir := filepath.Join(dataDir, "repos", "owner", "repo", "1", "default")
		Ok(t, os.MkdirAll(cloneDir, 0700))

		err := wd.DeletePlan(logger, safeRepo, safePull, "default", "../../../../etc", "")
		Assert(t, err != nil, "expected error for path traversal in project path")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})

	t.Run("MergeAgain rejects traversal in repo name", func(t *testing.T) {
		wdMerge := &events.FileWorkspace{
			DataDir:             dataDir,
			GpgNoSigningEnabled: true,
			CheckoutMerge:       true,
		}
		_, err := wdMerge.MergeAgain(logger, maliciousRepo, pull, "default")
		Assert(t, err != nil, "expected error for path traversal in repo name")
		Assert(t, strings.Contains(err.Error(), "traversal"), "expected traversal error, got: %s", err)
	})
}

func createPlanFile(t *testing.T, path string) {
	t.Helper()

// Copyright 2025 The Atlantis Authors
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrPathEscapesBase is returned when a path is outside the expected base directory.
var ErrPathEscapesBase = errors.New("path escapes base directory")

// EnsureSubPath returns an error if path is not contained within base.
// This prevents path traversal attacks where user-controlled data could
// escape the intended directory.
func EnsureSubPath(base, path string) error {
	cleanBase := filepath.Clean(base)
	cleanPath := filepath.Clean(path)
	// A path is within the base if it equals the base or starts with base + separator.
	if cleanPath != cleanBase && !strings.HasPrefix(cleanPath, cleanBase+string(os.PathSeparator)) {
		return ErrPathEscapesBase
	}
	return nil
}
// Copyright 2025 The Atlantis Authors
// SPDX-License-Identifier: Apache-2.0

package utils_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/runatlantis/atlantis/server/utils"
	. "github.com/runatlantis/atlantis/testing"
)

func TestEnsureSubPath(t *testing.T) {
	cases := []struct {
		base      string
		path      string
		expectErr bool
		desc      string
	}{
		{
			desc:      "path equals base",
			base:      "/data/repos",
			path:      "/data/repos",
			expectErr: false,
		},
		{
			desc:      "path within base",
			base:      "/data/repos",
			path:      "/data/repos/owner/repo/1/default",
			expectErr: false,
		},
		{
			desc:      "path traversal with .. staying within base",
			base:      "/data/repos",
			path:      "/data/repos/owner/repo/../../etc/passwd",
			expectErr: false, // resolves to /data/repos/etc/passwd which is within base
		},
		{
			desc:      "path traversal escaping base with ..",
			base:      "/data/repos",
			path:      "/data/repos/owner/../../../../etc/passwd",
			expectErr: true, // resolves to /etc/passwd which escapes base
		},
		{
			desc:      "path traversal escaping base",
			base:      "/data/repos",
			path:      "/etc/passwd",
			expectErr: true,
		},
		{
			desc:      "path with trailing separator in base",
			base:      "/data/repos/",
			path:      "/data/repos/owner/repo",
			expectErr: false,
		},
		{
			desc:      "base prefix but not subpath",
			base:      "/data/repos",
			path:      "/data/repos-evil/something",
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := utils.EnsureSubPath(tc.base, tc.path)
			if tc.expectErr {
				Assert(t, err != nil, "expected error but got nil")
				Assert(t, errors.Is(err, utils.ErrPathEscapesBase), "expected ErrPathEscapesBase but got %v", err)
				cleanBase := filepath.Clean(tc.base)
				cleanPath := filepath.Clean(tc.path)
				Assert(t, !strings.Contains(err.Error(), tc.base), "error should not expose base path: %q", err.Error())
				Assert(t, !strings.Contains(err.Error(), cleanBase), "error should not expose clean base path: %q", err.Error())
				Assert(t, !strings.Contains(err.Error(), tc.path), "error should not expose candidate path: %q", err.Error())
				Assert(t, !strings.Contains(err.Error(), cleanPath), "error should not expose clean candidate path: %q", err.Error())
			} else {
				Ok(t, err)
			}
		})
	}
}
