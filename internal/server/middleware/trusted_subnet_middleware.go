package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
)

func TrustedSubnetMiddleware(trustedSubnet string, logger *zap.Logger) gin.HandlerFunc {
	if trustedSubnet == "" {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	_, subnet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		logger.Fatal("invalid TrustedSubnet CIDR", zap.Error(err))
	}

	return func(c *gin.Context) {
		ipStr := c.GetHeader("X-Real-IP")
		if ipStr == "" {
			logger.Warn("missing X-Real-IP header")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			logger.Warn("invalid X-Real-IP header", zap.String("ip", ipStr))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !subnet.Contains(ip) {
			logger.Warn("IP not allowed", zap.String("ip", ipStr))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
