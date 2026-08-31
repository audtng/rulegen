package main

	if err != nil {
		return fmt.Errorf("decoding ciphertext: %w", err)
	}
	if len(fullCiphertext) < aesGCMNonceSize {
		return fmt.Errorf("ciphertext of length %d is too short, expected at least %d bytes", len(fullCiphertext), aesGCMNonceSize)
	}
	c.nonce = fullCiphertext[:aesGCMNonceSize]
	c.ciphertext = fullCiphertext[aesGCMNonceSize:]
	return nil
	}
	// decryptionRequest holds the request-specific ciphertextContainer and currently supported, optional query parameters: associatedData.
	decryptionRequest struct {
		CiphertextContainer *ciphertextContainer `json:"ciphertext"`
		AssociatedData      []byte               `json:"associated_data,omitempty"`
	}
	// encryptionResponse holds the response-specific ciphertextContainer.
	encryptionResponse struct {
			}, logger)
			return
		}
		if decReq.CiphertextContainer == nil {
			writeHTTPError(w, httpError{
				code:          http.StatusBadRequest,
				Errors:        []string{"missing mandatory field: ciphertext"},
				reqMethod:     r.Method,
				reqURI:        r.RequestURI,
				reqRemoteAddr: r.RemoteAddr,
			}, logger)
			return
		}
		key, err := deriveEncryptionKey(r.Context(), guard, fmt.Sprintf("%d_%s", decReq.CiphertextContainer.keyVersion, workloadSecretID))
		if err != nil {
			writeHTTPError(w, httpError{
			}, logger)
			return
		}
		plaintext, err := symmetricDecryptRaw(key, *decReq.CiphertextContainer, decReq.AssociatedData)
		if err != nil {
			writeHTTPError(w, httpError{
				code:          http.StatusBadRequest,
				Errors:        []string{fmt.Sprintf("decrypting: %v", err)},
				reqMethod:     r.Method,
				reqURI:        r.RequestURI,
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	})
}

func TestAdversarialDecryptInputs(t *testing.T) {
	testCases := map[string]struct {
		payload  string
		wantCode int
	}{
		"empty": {
			payload:  ``,
			wantCode: http.StatusBadRequest,
		},
		"no fields": {
			payload:  `{}`,
			wantCode: http.StatusBadRequest,
		},
		"truncated identifier": {
			payload:  `{"ciphertext":"vaul","associated_data":""}`,
			wantCode: http.StatusBadRequest,
		},
		"truncated nonce": {
			payload:  `{"ciphertext":"vault:v1:AAAA","associated_data":""}`,
			wantCode: http.StatusBadRequest,
		},
		"random data": {
			payload:  `{"ciphertext":"vault:v1:DEOLYTbWcdjNqFjrfb/heTN7P1LohVS8KpGYisecIs2O2FrjevU3zrJHPTe6biaalMa2xphSNt6JFvWpCSsW4svHe2r0myjnq05Zbi+h37GIS2FYPrRDnoglsRHdJ+rH+D0MBJON0WQnVE9qbSDI4P9cjZZXVm7lx2VAqu3ioWw=","associated_data":""}`,
			wantCode: http.StatusBadRequest,
		},
		//
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			fakeStateGuard, err := newTestGuard()
			require.NoError(err)
			mux := newMockTransitEngineMux(fakeStateGuard)

			decryptReq := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/v1/transit/decrypt/foo", bytes.NewReader([]byte(tc.payload)))
			decryptReq.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, decryptReq)
			res := rec.Result()
			t.Cleanup(func() { _ = res.Body.Close() })
			body, err := io.ReadAll(res.Body)
			require.NoError(err)
			require.Equal(tc.wantCode, res.StatusCode, string(body))
		})
	}
}

type fakeStateGuard struct {
	state *stateguard.State
}
