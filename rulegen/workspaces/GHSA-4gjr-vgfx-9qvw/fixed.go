package main

package middlewares

import (
	"net/url"
	stdpath "path"

	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func FsUp(c *gin.Context) {
	path := c.GetHeader("File-Path")
	password := c.GetHeader("Password")
	path, err := url.PathUnescape(path)
	if err != nil {
		common.ErrorResp(c, err, 400)
		c.Abort()
		return
	}
	user := c.MustGet("user").(*model.User)
	path = stdpath.Join(user.BasePath, path)
	meta, err := db.GetNearestMeta(stdpath.Dir(path))
	if err != nil {
		if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
			common.ErrorResp(c, err, 500, true)
			c.Abort()
			return
		}
	}
	if !(common.CanAccess(user, meta, path, password) && (user.CanWrite() || common.CanWrite(meta, path))) {
		common.ErrorResp(c, errs.PermissionDenied, 403)
		c.Abort()
		return
	}
	c.Next()
}
package common

import (
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
)

func CanWrite(meta *model.Meta, path string) bool {
	if meta == nil || !meta.Write {
		return false
	}
	return meta.WSub || meta.Path == path
}

func CanAccess(user *model.User, meta *model.Meta, path string, password string) bool {
	// if is not guest, can access
	if user.CanAccessWithoutPassword() {
		return true
	}
	// if meta is nil or password is empty, can access
	if meta == nil || meta.Password == "" {
		return true
	}
	// if meta doesn't apply to sub_folder, can access
	if !utils.PathEqual(meta.Path, path) && !meta.PSub {
		return true
	}
	// validate password
	return meta.Password == password
}
	"strconv"
	"time"

	"github.com/alist-org/alist/v3/internal/fs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/gin-gonic/gin"
)

func FsStream(c *gin.Context) {
	asTask := c.GetHeader("As-Task") == "true"
	user := c.MustGet("user").(*model.User)
	path = stdpath.Join(user.BasePath, path)

	dir, name := stdpath.Split(path)
	sizeStr := c.GetHeader("Content-Length")
	asTask := c.GetHeader("As-Task") == "true"
	user := c.MustGet("user").(*model.User)
	path = stdpath.Join(user.BasePath, path)

	storage, err := fs.GetStorage(path)
	if err != nil {
		common.ErrorResp(c, err, 400)
