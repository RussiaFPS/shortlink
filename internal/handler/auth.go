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

// encryptCookie encrypts a cookie value.
func encryptCookie(cookieValue string, secretKey []byte) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(cookieValue))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(cookieValue))
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decryptCookie decrypts a cookie value.
func decryptCookie(cipherText string, secretKey []byte) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	initVector := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, initVector)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
