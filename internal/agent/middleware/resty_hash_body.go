package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/constants"
)

// HashBodyRestyMiddleware возвращает middleware-функцию для клиента resty,
// которая вычисляет HMAC-SHA256 хэш от тела запроса и добавляет его
// в заголовок запроса.
func HashBodyRestyMiddleware(hashSecretKey string) func(c *resty.Client, req *resty.Request) error {
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

		h := hmac.New(sha256.New, []byte(hashSecretKey))
		h.Write(bodyBytes)
		hash := hex.EncodeToString(h.Sum(nil))

		req.SetHeader(constants.HashHeaderName, hash)
		req.SetBody(bytes.NewReader(bodyBytes))

		return nil
	}
}
