package main

}

func (err ErrFileTypeForbidden) Error() string {
	return "This file extension or type is not allowed to be uploaded."
}

var wildcardTypeRe = regexp.MustCompile(`^[a-z]+/\*$`)

// Verify validates whether a file is allowed to be uploaded.
func Verify(buf []byte, fileName, allowedTypesStr string) error {
	allowedTypesStr = strings.ReplaceAll(allowedTypesStr, "|", ",") // compat for old config format

		return ErrFileTypeForbidden{Type: fullMimeType}
	}
	extension := strings.ToLower(path.Ext(fileName))

	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#Unique_file_type_specifiers
	for _, allowEntry := range allowedTypes {
		if allowEntry == "*/*" {
			return nil // everything allowed
		} else if strings.HasPrefix(allowEntry, ".") && allowEntry == extension {
			return nil // extension is allowed
		} else if mimeType == allowEntry {
			return nil // mime type is allowed
		} else if wildcardTypeRe.MatchString(allowEntry) && strings.HasPrefix(mimeType, allowEntry[:len(allowEntry)-1]) {
			return nil // wildcard match, e.g. image/*
		}
	}

	log.Info("Attachment with type %s blocked from upload", fullMimeType)
	return ErrFileTypeForbidden{Type: fullMimeType}
}

	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/attachment"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/context/upload"
	"code.gitea.io/gitea/services/convert"
		filename = query
	}

	attachment, err := attachment.UploadAttachment(ctx, file, setting.Attachment.AllowedTypes, header.Size, &repo_model.Attachment{
		Name:       filename,
		UploaderID: ctx.Doer.ID,
		RepoID:     ctx.Repo.Repository.ID,
	//     "$ref": "#/responses/Attachment"
	//   "404":
	//     "$ref": "#/responses/error"
	//   "423":
	//     "$ref": "#/responses/repoArchivedError"

		attachment.Name = form.Name
	}

	if err := repo_model.UpdateAttachment(ctx, attachment); err != nil {
		ctx.Error(http.StatusInternalServerError, "UpdateAttachment", err)
	}

	ctx.JSON(http.StatusCreated, convert.ToAPIAttachment(ctx.Repo.Repository, attachment))
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/attachment"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/context/upload"
	"code.gitea.io/gitea/services/convert"
		filename = query
	}

	attachment, err := attachment.UploadAttachment(ctx, file, setting.Attachment.AllowedTypes, header.Size, &repo_model.Attachment{
		Name:       filename,
		UploaderID: ctx.Doer.ID,
		RepoID:     ctx.Repo.Repository.ID,
	//     "$ref": "#/responses/Attachment"
	//   "404":
	//     "$ref": "#/responses/error"
	//   "423":
	//     "$ref": "#/responses/repoArchivedError"
	attach := getIssueCommentAttachmentSafeWrite(ctx)
		attach.Name = form.Name
	}

	if err := repo_model.UpdateAttachment(ctx, attach); err != nil {
		ctx.Error(http.StatusInternalServerError, "UpdateAttachment", attach)
	}
	ctx.JSON(http.StatusCreated, convert.ToAPIAttachment(ctx.Repo.Repository, attach))
}
