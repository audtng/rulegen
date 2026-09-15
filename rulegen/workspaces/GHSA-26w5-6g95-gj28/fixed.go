package main


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
