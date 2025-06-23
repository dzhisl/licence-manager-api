package middleware

import (
	"context"
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"github.com/spf13/viper"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func AdminAuthMiddleware(c *gin.Context) {
	adminKey := viper.GetString("ADMIN_SECRET_KEY")
	if adminKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		c.Abort()
		return
	}

	header := c.GetHeader("X-API-Key")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
		c.Abort()
		return
	}

	if subtle.ConstantTimeCompare([]byte(header), []byte(adminKey)) == 1 {
		c.Next()
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
	c.Abort()
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
