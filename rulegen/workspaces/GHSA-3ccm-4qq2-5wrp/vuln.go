package main

	if err != nil {
		return fmt.Errorf("decoding ciphertext: %w", err)
	}
	c.nonce = fullCiphertext[:aesGCMNonceSize]
	c.ciphertext = fullCiphertext[aesGCMNonceSize:]
	return nil
	}
	// decryptionRequest holds the request-specific ciphertextContainer and currently supported, optional query parameters: associatedData.
	decryptionRequest struct {
		CiphertextContainer ciphertextContainer `json:"ciphertext"`
		AssociatedData      []byte              `json:"associated_data,omitempty"`
	}
	// encryptionResponse holds the response-specific ciphertextContainer.
	encryptionResponse struct {
			}, logger)
			return
		}
		key, err := deriveEncryptionKey(r.Context(), guard, fmt.Sprintf("%d_%s", decReq.CiphertextContainer.keyVersion, workloadSecretID))
		if err != nil {
			writeHTTPError(w, httpError{
			}, logger)
			return
		}
		plaintext, err := symmetricDecryptRaw(key, decReq.CiphertextContainer, decReq.AssociatedData)
		if err != nil {
			writeHTTPError(w, httpError{
				code:          http.StatusInternalServerError,
				Errors:        []string{fmt.Sprintf("decrypting: %v", err)},
				reqMethod:     r.Method,
				reqURI:        r.RequestURI,
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	})
}

type fakeStateGuard struct {
	state *stateguard.State
}
