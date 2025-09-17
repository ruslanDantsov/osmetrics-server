package middleware

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/crypto"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func NewDecryptPayloadMiddleware(privKey *rsa.PrivateKey, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if privKey == nil {
			c.Next()
			return
		}

		// читаем тело
		encryptedBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("failed to read request body", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		defer func() {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(encryptedBody))
		}()

		// сначала декодируем base64
		cipherData, err := base64.StdEncoding.DecodeString(string(encryptedBody))
		if err != nil {
			logger.Error("failed to decode base64 payload", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// расшифровываем RSA
		plain, err := crypto.DecryptRSA(privKey, cipherData)
		if err != nil {
			logger.Error("failed to decrypt payload", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// подменяем тело на JSON(Metrics)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(plain))
		c.Next()
	}
}
