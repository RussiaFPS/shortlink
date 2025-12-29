package handler

import (
	"github.com/RussiaFPS/shortlink/internal/service"
	"github.com/gin-gonic/gin"
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
			encryptedCookie, uid, expiration, err = service.GenCookie(secretKey)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie("auth", encryptedCookie, int(expiration.Unix()), "/", "", false, true)
			c.Set(string(UserIDKey), uid)
			c.Next()
			return
		}

		uid, err = service.DecryptCookie(cookie, []byte(*secretKey))
		if err != nil {
			encryptedCookie, uid, expiration, err = service.GenCookie(secretKey)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie("auth", encryptedCookie, int(expiration.Unix()), "/", "", false, true)
		}
		c.Set(string(UserIDKey), uid)
		c.Next()
	}
}
