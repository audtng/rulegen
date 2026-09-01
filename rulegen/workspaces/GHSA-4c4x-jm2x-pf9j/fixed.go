package main

			hasher := sha256.New()
			var tee io.Reader
			if isURL(artifactStr) {
				r, err := util.FileOrURLReadCloser(cmd.Context(), artifactStr, nil)
				if err != nil {
					return nil, fmt.Errorf("error fetching '%v': %w", artifactStr, err)
				}
				defer r.Close()
				tee = io.TeeReader(r, hasher)
			} else {
				file, err := os.Open(filepath.Clean(artifactStr))
				if err != nil {
			splitPubKeyString := strings.Split(publicKeyStr, ",")
			if len(splitPubKeyString) == 1 {
				if isURL(splitPubKeyString[0]) {
					r, err := util.FileOrURLReadCloser(cmd.Context(), splitPubKeyString[0], nil)
					if err != nil {
						return nil, fmt.Errorf("error fetching '%v': %w", splitPubKeyString[0], err)
					}
					defer r.Close()
					c, err := io.ReadAll(r)
					if err != nil {
						return nil, fmt.Errorf("error reading public key from '%v': %w", splitPubKeyString[0], err)
					}
					params.Query.PublicKey.Content = c
				} else {
					keyBytes, err := os.ReadFile(filepath.Clean(splitPubKeyString[0]))
					if err != nil {
	"fmt"
	"io"
	"net/http"
	"time"
)

// FileOrURLReadCloser reads content either from a URL or a byte slice
// Note: Caller is responsible for closing the returned ReadCloser
// Note: This must never be called from any server codepath to prevent SSRF
func FileOrURLReadCloser(ctx context.Context, url string, content []byte) (io.ReadCloser, error) {
	var dataReader io.ReadCloser
	if url != "" {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
		if err != nil {
			return handleRekorAPIError(params, http.StatusBadRequest, err, unsupportedPKIFormat)
		}
		keyReader := bytes.NewReader(params.Query.PublicKey.Content)

		key, err := af.NewPublicKey(keyReader)
		if err != nil {
