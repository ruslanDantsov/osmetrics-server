package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// NewGzipDecompressionMiddleware возвращает middleware для Gin,
// который автоматически распаковывает тело запроса, если оно
// сжато с помощью gzip.
func NewGzipDecompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := c.GetHeader("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			gzipReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Invalid gzip data",
				})
				return
			}
			defer func() {
				if err := gzipReader.Close(); err != nil {
					err = c.Error(err)
				}
			}()

			// Replace the request body with the decompressed stream
			c.Request.Body = io.NopCloser(gzipReader)
		}

		c.Next()
	}
}
