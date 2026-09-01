package main

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"net/http"
	"net/url"
)

// Service providing Discord as a notification service
)

// Send a notification message to discord
func (service *Service) Send(message string, params *types.Params) error {

	if service.config.JSON {
		postURL := CreateAPIURLFromConfig(service.config)
		return doSend([]byte(message), postURL)
	}

	items, omitted := CreateItemsFromPlain(message, service.config.SplitLines)
	return service.sendItems(items, params, omitted)
}

// SendItems sends items with additional meta data and richer appearance
		err = fmt.Errorf("response status code %s", res.Status)
	}

	if err != nil {
		return fmt.Errorf("failed to send discord notification: %v", err)
	}

	return nil
}

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	"time"
)

// WebhookPayload is the webhook endpoint payload
type WebhookPayload struct {
	Embeds   []embedItem `json:"embeds"`
	Username string      `json:"username,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
}

// JSON is the actual notification payload
// CreatePayloadFromItems creates a JSON payload to be sent to the discord webhook API
func CreatePayloadFromItems(items []types.MessageItem, title string, colors [types.MessageLevelCount]uint, omitted int) (WebhookPayload, error) {

	metaCount := 1
	if omitted < 1 && len(title) < 1 {
		metaCount = 0
		embeds = append(embeds, ei)
	}

	embeds[0].Title = title
	if omitted > 0 {
		embeds[0].Footer = &embedFooter{
			Text: fmt.Sprintf("... (%v character(s) where omitted)", omitted),
		}
	}

	maxTotal := Min(len(runes), limits.TotalChunkSize)
	maxCount := limits.ChunkCount - 1

	for i := 0; i < maxCount; i++ {
		// If no suitable split point is found, use the chunkSize
		chunkEnd := chunkOffset + limits.ChunkSize
		// ... and start next chunk directly after this one
		nextChunkStart := chunkEnd
		if chunkEnd > maxTotal {
			// The chunk is smaller than the limit, no need to search
			chunkEnd = maxTotal
