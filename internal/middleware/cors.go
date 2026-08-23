package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	isDev := os.Getenv("APP_ENV") == "development"
	allowedOriginsRaw := os.Getenv("CORS_ALLOWED_ORIGINS")

	allowedMap := make(map[string]bool)
	for _, origin := range strings.Split(allowedOriginsRaw, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowedMap[trimmed] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowOrigin := ""
		if isDev {
			if origin != "" {
				allowOrigin = origin
			} else {
				allowOrigin = "*"
			}
		} else if origin != "" && allowedMap[origin] {
			allowOrigin = origin
		}

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			if isDev || allowOrigin != "" {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		c.Next()
	}
}
