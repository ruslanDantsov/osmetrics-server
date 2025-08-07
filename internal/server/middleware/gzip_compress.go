// Package middleware provides HTTP middleware handlers for the application.
package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

func NewGzipCompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		gz := gzip.NewWriter(c.Writer)
		defer func() {
			if err := gz.Close(); err != nil {
				err = c.Error(err)
			}
		}()

		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			Writer:         gz,
		}

		c.Header("Content-Encoding", "gzip")

		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Writer.Write([]byte(s))
}
