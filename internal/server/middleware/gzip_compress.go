package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

// NewGzipCompressionMiddleware возвращает middleware для сжатия HTTP-ответов в формате gzip.
func NewGzipCompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.Next()
			return
		}
		defer gz.Close()

		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			Writer:         gz,
		}

		c.Header("Content-Encoding", "gzip")

		c.Next()
	}
}

// gzipResponseWriter оборачивает gin.ResponseWriter и реализует запись
// данных в сжатом формате gzip через io.Writer.
type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write записывает сжатые байты данных в поток ответа.
func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

// WriteString записывает сжатую строку в поток ответа.
func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Writer.Write([]byte(s))
}
