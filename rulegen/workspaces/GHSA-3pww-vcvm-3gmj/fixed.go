package main

	}

	// since user only get notifications when he has access to use minimal access mode
	if n.Repository == nil {
		return result
	}
	perm, err := access_model.GetIndividualUserRepoPermission(ctx, n.Repository, n.User)
	if err != nil {
		log.Error("GetIndividualUserRepoPermission failed: %v", err)
		return result
	}
	// if the user has been revoked access to the repo, do not leak repo or subject info
	if !perm.HasAnyUnitAccessOrPublicAccess() {
		return result
	}
	result.Repository = ToRepo(ctx, n.Repository, perm)
	// This permission is not correct and we should not be reporting it
	for repository := result.Repository; repository != nil; repository = repository.Parent {
		repository.Permissions = nil
	}

	// handle Subject
package feed

import (
	auth_model "gitea.dev/models/auth"
	"gitea.dev/services/context"
)

// checkRepoFeedTokenScope ensures an API token has repository read scope before a
// feed serves private repository content, mirroring checkDownloadTokenScope for
// downloads. Returns false (and writes the response) when the token is denied.
func checkRepoFeedTokenScope(ctx *context.Context) bool {
	context.CheckRepoScopedToken(ctx, ctx.Repo.Repository, auth_model.Read)
	return !ctx.Written()
}

// RenderBranchFeed render format for branch or file
func RenderBranchFeed(ctx *context.Context, feedType string) {
	if ctx.Repo.TreePath == "" {
	}

	cmd := gitcmd.NewCommand().AddArguments("clone")
	// Never follow HTTP redirects: no clone caller needs them, and a remote redirecting to an
	// otherwise-blocked address would be an SSRF vector (e.g. migrating from an attacker URL).
	cmd.AddArguments("-c", "http.followRedirects=false")
	if opts.SkipTLSVerify {
		cmd.AddArguments("-c", "http.sslVerify=false")
	}
