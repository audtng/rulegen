package main

package repo

import (
	"net/http"

	issues_model "code.gitea.io/gitea/models/issues"
	access_model "code.gitea.io/gitea/models/perm/access"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unit"
	"code.gitea.io/gitea/modules/httpcache"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"

	file, header, err := ctx.Req.FormFile("file")
	if err != nil {
		ctx.ServerError("FormFile", err)
		return
	}
	defer file.Close()
			ctx.HTTPError(http.StatusBadRequest, err.Error())
			return
		}
		ctx.ServerError("UploadAttachmentGeneralSizeLimit", err)
		return
	}

		ctx.HTTPError(http.StatusBadRequest, err.Error())
		return
	}

	if !ctx.IsSigned {
		ctx.HTTPError(http.StatusForbidden)
		return
	}

	if attach.RepoID != ctx.Repo.Repository.ID {
		ctx.HTTPError(http.StatusBadRequest, "attachment does not belong to this repository")
		return
	}

	if ctx.Doer.ID != attach.UploaderID {
		if attach.IssueID > 0 {
			issue, err := issues_model.GetIssueByID(ctx, attach.IssueID)
			if err != nil {
				ctx.ServerError("GetIssueByID", err)
				return
			}
			if !ctx.Repo.Permission.CanWriteIssuesOrPulls(issue.IsPull) {
				ctx.HTTPError(http.StatusForbidden)
				return
			}
		} else if attach.ReleaseID > 0 {
			if !ctx.Repo.Permission.CanWrite(unit.TypeReleases) {
				ctx.HTTPError(http.StatusForbidden)
				return
			}
		} else {
			if !ctx.Repo.Permission.IsAdmin() && !ctx.Repo.Permission.IsOwner() {
				ctx.HTTPError(http.StatusForbidden)
				return
			}
		}
	}

	err = repo_model.DeleteAttachment(ctx, attach, true)
	if err != nil {
		ctx.ServerError("DeleteAttachment", err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
	} else { // If we have the repository we check access
		perm, err := access_model.GetUserRepoPermission(ctx, repository, ctx.Doer)
		if err != nil {
			ctx.ServerError("GetUserRepoPermission", err)
			return
		}
		if !perm.CanRead(unitType) {
		ctx.APIErrorNotFound()
		return
	}

	if err := repo_model.DeleteAttachment(ctx, attach, true); err != nil {
		ctx.APIErrorInternal(err)
