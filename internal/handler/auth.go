package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

// AuthenticationMiddleware is a middleware that handles user authentication.
func AuthenticationMiddleware(secretKey *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		var encryptedCookie, uid string
		var expiration *time.Time

		if err != nil {
			encryptedCookie, uid, expiration, err = genCookie(secretKey)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie("auth", encryptedCookie, int(expiration.Unix()), "/", "", false, true)
			c.Set("uid", uid)
			c.Next()
			return
		}

		uid, err = decryptCookie(cookie, []byte(*secretKey))
		if err != nil {
			encryptedCookie, uid, expiration, err = genCookie(secretKey)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie("auth", encryptedCookie, int(expiration.Unix()), "/", "", false, true)
		}
		c.Set("uid", uid)
		c.Next()
	}
}

// genCookie generates a new cookie.
func genCookie(secretKey *string) (string, string, *time.Time, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", "", nil, err
	}
	uid := newUUID.String()
	encryptedCookie, err := encryptCookie(uid, []byte(*secretKey))
	if err != nil {
		return "", "", nil, err
	}
	expiration := time.Now().Add(24 * time.Hour)

	return encryptedCookie, uid, &expiration, nil
}

// encryptCookie encrypts a cookie value using AES-GCM.
func encryptCookie(cookieValue string, secretKey []byte) (string, error) {
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

// decryptCookie decrypts a cookie value using AES-GCM.
func decryptCookie(cipherText string, secretKey []byte) (string, error) {
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
