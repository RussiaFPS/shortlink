package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

// GenCookie generates a new cookie.
func GenCookie(secretKey *string) (string, string, *time.Time, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", "", nil, err
	}
	uid := newUUID.String()
	encryptedCookie, err := EncryptCookie(uid, []byte(*secretKey))
	if err != nil {
		return "", "", nil, err
	}
	expiration := time.Now().Add(24 * time.Hour)

	return encryptedCookie, uid, &expiration, nil
}

// EncryptCookie encrypts a cookie value using AES-GCM.
func EncryptCookie(cookieValue string, secretKey []byte) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(cookieValue), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptCookie decrypts a cookie value using AES-GCM.
func DecryptCookie(cipherText string, secretKey []byte) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
