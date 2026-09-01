package main

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// Service providing Discord as a notification service
)

// Send a notification message to discord
func (service *Service) Send(message string, params *types.Params) (err error) {
	if service.config.JSON {
		postURL := CreateAPIURLFromConfig(service.config)
		err = doSend([]byte(message), postURL)
	} else {
		items, omitted := CreateItemsFromPlain(message, service.config.SplitLines)
		err = service.sendItems(items, params, omitted)
	}

	if err != nil {
		err = fmt.Errorf("failed to send discord notification: %v", err)
	}

	return
}

// SendItems sends items with additional meta data and richer appearance
		err = fmt.Errorf("response status code %s", res.Status)
	}

	return err
}

import (
	"fmt"
	"time"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// WebhookPayload is the webhook endpoint payload
type WebhookPayload struct {
	Embeds    []embedItem `json:"embeds"`
	Username  string      `json:"username,omitempty"`
	AvatarURL string      `json:"avatar_url,omitempty"`
}

// JSON is the actual notification payload
// CreatePayloadFromItems creates a JSON payload to be sent to the discord webhook API
func CreatePayloadFromItems(items []types.MessageItem, title string, colors [types.MessageLevelCount]uint, omitted int) (WebhookPayload, error) {

	if len(items) < 1 {
		return WebhookPayload{}, fmt.Errorf("message is empty")
	}

	metaCount := 1
	if omitted < 1 && len(title) < 1 {
		metaCount = 0
		embeds = append(embeds, ei)
	}

	// This should not happen, but it's better to leave the index check before dereferencing the array
	if len(embeds) > 0 {
		embeds[0].Title = title

		if omitted > 0 {
			embeds[0].Footer = &embedFooter{
				Text: fmt.Sprintf("... (%v character(s) were omitted)", omitted),
			}
		}
	}

	maxTotal := Min(len(runes), limits.TotalChunkSize)
	maxCount := limits.ChunkCount - 1

	if len(input) == 0 {
		// If the message is empty, return an empty array
		omitted = 0
		return
	}

	for i := 0; i < maxCount; i++ {
		// If no suitable split point is found, start next chunk at chunkSize from chunk start
		nextChunkStart := chunkOffset + limits.ChunkSize
		// ... and set the chunk end to the rune before the next chunk
		chunkEnd := nextChunkStart - 1
		if chunkEnd > maxTotal {
			// The chunk is smaller than the limit, no need to search
			chunkEnd = maxTotal
