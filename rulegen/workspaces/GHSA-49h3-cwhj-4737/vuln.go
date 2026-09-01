package main

			return fmt.Errorf("failed to decode X-WOPI-SuggestedTarget header (UTF-7): %w", err)
		}

		fileUriParsed, err := fs.NewUriFromString(fileUri)
		if err != nil {
			return fmt.Errorf("failed to parse file uri: %w", err)
	return nil
}

func (service *WopiService) GetFile(c *gin.Context) error {
	uri, m, _, viewerSession, dep, err := prepareFs(c)
	if err != nil {

type (
	CreateViewerSessionService struct {
		Uri             string               `json:"uri" form:"uri" binding:"required"`
		Version         string               `json:"version" form:"version"`
		ViewerID        string               `json:"viewer_id" form:"viewer_id" binding:"required"`
		PreferredAction types.ViewerAction `json:"preferred_action" form:"preferred_action" binding:"required"`
	}
	CreateViewerSessionParamCtx struct{}
	"context"
	"errors"
	"fmt"
	"github.com/cloudreve/Cloudreve/v4/inventory/types"
	"net/url"
	"strings"
	"time"

	"github.com/cloudreve/Cloudreve/v4/application/dependency"
	"github.com/cloudreve/Cloudreve/v4/pkg/cluster/routes"
	"github.com/cloudreve/Cloudreve/v4/pkg/filemanager/manager"
	LockTokenHeader       = WopiHeaderPrefix + "Lock"
	ItemVersionHeader     = WopiHeaderPrefix + "ItemVersion"
	SuggestedTargetHeader = WopiHeaderPrefix + "SuggestedTarget"

	MethodLock           = "LOCK"
	MethodUnlock         = "UNLOCK"
