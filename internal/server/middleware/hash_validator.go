package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func HashCheckerMiddleware(hashSecretKey string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" && c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		agentHash := c.GetHeader(constants.HashHeaderName)
		if len(agentHash) > 0 {

			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error("Failed to read request body",
					zap.Error(err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			h := hmac.New(sha256.New, []byte(hashSecretKey))
			h.Write(body)
			serverHash := hex.EncodeToString(h.Sum(nil))

			if !hmac.Equal([]byte(serverHash), []byte(agentHash)) {
				logger.Error("Invalid hash signature",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("expected", serverHash),
					zap.String("received", agentHash),
				)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request signature",
				})
				return
			}

			c.Header(constants.HashHeaderName, serverHash)
		}

		c.Next()
	}
}
