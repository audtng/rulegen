package main

	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)
}

func (gcm *AesGCM) Decrypt(cipherText, nonce []byte) ([]byte, error) {
	if len(nonce) != gcm.NonceSize() {
		return nil, errors.New("incorrect nonce length")
	}

	plainText, err := gcm.Open(nil, nonce, cipherText, []byte{})
	if err != nil {
		return nil, err
