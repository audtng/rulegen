package main

	"strconv"
	"time"

	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/fs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func FsStream(c *gin.Context) {
	asTask := c.GetHeader("As-Task") == "true"
	user := c.MustGet("user").(*model.User)
	path = stdpath.Join(user.BasePath, path)
	if !user.CanWrite() {
		meta, err := db.GetNearestMeta(stdpath.Dir(path))
		if err != nil {
			if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
				common.ErrorResp(c, err, 500, true)
				return
			}
		}
		if !canWrite(meta, path) {
			common.ErrorResp(c, errs.PermissionDenied, 403)
			return
		}
	}

	dir, name := stdpath.Split(path)
	sizeStr := c.GetHeader("Content-Length")
	asTask := c.GetHeader("As-Task") == "true"
	user := c.MustGet("user").(*model.User)
	path = stdpath.Join(user.BasePath, path)
	if !user.CanWrite() {
		meta, err := db.GetNearestMeta(stdpath.Dir(path))
		if err != nil {
			if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
				common.ErrorResp(c, err, 500, true)
				return
			}
		}
		if !canWrite(meta, path) {
			common.ErrorResp(c, errs.PermissionDenied, 403)
			return
		}
	}
	storage, err := fs.GetStorage(path)
	if err != nil {
		common.ErrorResp(c, err, 400)
