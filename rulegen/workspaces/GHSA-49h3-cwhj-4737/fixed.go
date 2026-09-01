package main

			return fmt.Errorf("failed to decode X-WOPI-SuggestedTarget header (UTF-7): %w", err)
		}

		// X-WOPI-SuggestedTarget is a filename, not a path. Reject any value
		// that could traverse out of the source file's directory.
		if !isValidWopiSuggestedTarget(fileName) {
			c.Status(http.StatusBadRequest)
			c.Header(wopi.InvalidFileNameHeader, "Invalid target file name")
			return nil
		}

		fileUriParsed, err := fs.NewUriFromString(fileUri)
		if err != nil {
			return fmt.Errorf("failed to parse file uri: %w", err)
	return nil
}

// isValidWopiSuggestedTarget enforces that X-WOPI-SuggestedTarget is a bare
// filename (or extension prefixed with ".") and cannot traverse out of the
// source file's directory once joined onto it.
func isValidWopiSuggestedTarget(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	if strings.ContainsAny(name, "/\\") {
		return false
	}
	return true
}

func (service *WopiService) GetFile(c *gin.Context) error {
	uri, m, _, viewerSession, dep, err := prepareFs(c)
	if err != nil {

type (
	CreateViewerSessionService struct {
		Uri             string             `json:"uri" form:"uri" binding:"required"`
		Version         string             `json:"version" form:"version"`
		ViewerID        string             `json:"viewer_id" form:"viewer_id" binding:"required"`
		PreferredAction types.ViewerAction `json:"preferred_action" form:"preferred_action" binding:"required"`
	}
	CreateViewerSessionParamCtx struct{}
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cloudreve/Cloudreve/v4/inventory/types"

	"github.com/cloudreve/Cloudreve/v4/application/dependency"
	"github.com/cloudreve/Cloudreve/v4/pkg/cluster/routes"
	"github.com/cloudreve/Cloudreve/v4/pkg/filemanager/manager"
	LockTokenHeader       = WopiHeaderPrefix + "Lock"
	ItemVersionHeader     = WopiHeaderPrefix + "ItemVersion"
	SuggestedTargetHeader = WopiHeaderPrefix + "SuggestedTarget"
	InvalidFileNameHeader = WopiHeaderPrefix + "InvalidFileNameError"

	MethodLock           = "LOCK"
	MethodUnlock         = "UNLOCK"
