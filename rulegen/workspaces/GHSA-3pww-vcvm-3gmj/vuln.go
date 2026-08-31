package main

	}

	// since user only get notifications when he has access to use minimal access mode
	if n.Repository != nil {
		perm, err := access_model.GetIndividualUserRepoPermission(ctx, n.Repository, n.User)
		if err != nil {
			log.Error("GetIndividualUserRepoPermission failed: %v", err)
			return result
		}
		if perm.HasAnyUnitAccessOrPublicAccess() { // if user has been revoked access to repo, do not show repo info
			result.Repository = ToRepo(ctx, n.Repository, perm)
			// This permission is not correct and we should not be reporting it
			for repository := result.Repository; repository != nil; repository = repository.Parent {
				repository.Permissions = nil
			}
		}
	}

	// handle Subject
package feed

import (
	"gitea.dev/services/context"
)

// RenderBranchFeed render format for branch or file
func RenderBranchFeed(ctx *context.Context, feedType string) {
	if ctx.Repo.TreePath == "" {
	}

	cmd := gitcmd.NewCommand().AddArguments("clone")
	if opts.SkipTLSVerify {
		cmd.AddArguments("-c", "http.sslVerify=false")
	}
