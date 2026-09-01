package main

	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)
}

func (gcm *AesGCM) Decrypt(cipherText, nonce []byte) ([]byte, error) {
	plainText, err := gcm.Open(nil, nonce, cipherText, []byte{})
	if err != nil {
		return nil, err
