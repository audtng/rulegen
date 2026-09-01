package main

package repo

import (
	"fmt"
	"net/http"

	access_model "code.gitea.io/gitea/models/perm/access"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/modules/httpcache"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"

	file, header, err := ctx.Req.FormFile("file")
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError, fmt.Sprintf("FormFile: %v", err))
		return
	}
	defer file.Close()
			ctx.HTTPError(http.StatusBadRequest, err.Error())
			return
		}
		ctx.HTTPError(http.StatusInternalServerError, fmt.Sprintf("NewAttachment: %v", err))
		return
	}

		ctx.HTTPError(http.StatusBadRequest, err.Error())
		return
	}
	if !ctx.IsSigned || (ctx.Doer.ID != attach.UploaderID) {
		ctx.HTTPError(http.StatusForbidden)
		return
	}
	err = repo_model.DeleteAttachment(ctx, attach, true)
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError, fmt.Sprintf("DeleteAttachment: %v", err))
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
	} else { // If we have the repository we check access
		perm, err := access_model.GetUserRepoPermission(ctx, repository, ctx.Doer)
		if err != nil {
			ctx.HTTPError(http.StatusInternalServerError, "GetUserRepoPermission", err.Error())
			return
		}
		if !perm.CanRead(unitType) {
		ctx.APIErrorNotFound()
		return
	}
	// FIXME Should prove the existence of the given repo, but results in unnecessary database requests

	if err := repo_model.DeleteAttachment(ctx, attach, true); err != nil {
		ctx.APIErrorInternal(err)
