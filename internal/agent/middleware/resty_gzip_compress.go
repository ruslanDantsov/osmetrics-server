package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/go-resty/resty/v2"
)

func GzipRestyMiddleware() func(c *resty.Client, req *resty.Request) error {
	return func(c *resty.Client, req *resty.Request) error {
		if req.Body == nil {
			return nil
		}

		var bodyBytes []byte
		switch v := req.Body.(type) {
		case []byte:
			bodyBytes = v
		case string:
			bodyBytes = []byte(v)
		default:
			return nil
		}

		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		if _, err := gz.Write(bodyBytes); err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}

		req.SetBody(buf.Bytes())
		req.SetHeader("Content-Encoding", "gzip")

		return nil
	}
}
