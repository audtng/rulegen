package main

			hasher := sha256.New()
			var tee io.Reader
			if isURL(artifactStr) {
				/* #nosec G107 */
				resp, err := http.Get(artifactStr)
				if err != nil {
					return nil, fmt.Errorf("error fetching '%v': %w", artifactStr, err)
				}
				defer resp.Body.Close()
				tee = io.TeeReader(resp.Body, hasher)
			} else {
				file, err := os.Open(filepath.Clean(artifactStr))
				if err != nil {
			splitPubKeyString := strings.Split(publicKeyStr, ",")
			if len(splitPubKeyString) == 1 {
				if isURL(splitPubKeyString[0]) {
					params.Query.PublicKey.URL = strfmt.URI(splitPubKeyString[0])
				} else {
					keyBytes, err := os.ReadFile(filepath.Clean(splitPubKeyString[0]))
					if err != nil {
	"fmt"
	"io"
	"net/http"
)

// FileOrURLReadCloser Note: caller is responsible for closing ReadCloser returned from method!
func FileOrURLReadCloser(ctx context.Context, url string, content []byte) (io.ReadCloser, error) {
	var dataReader io.ReadCloser
	if url != "" {
		//TODO: set timeout here, SSL settings?
		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
		if err != nil {
			return handleRekorAPIError(params, http.StatusBadRequest, err, unsupportedPKIFormat)
		}
		keyReader, err := util.FileOrURLReadCloser(httpReqCtx, params.Query.PublicKey.URL.String(), params.Query.PublicKey.Content)
		if err != nil {
			return handleRekorAPIError(params, http.StatusBadRequest, err, malformedPublicKey)
		}
		defer keyReader.Close()

		key, err := af.NewPublicKey(keyReader)
		if err != nil {
