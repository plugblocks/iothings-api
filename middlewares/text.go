package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/services"
)

func TextMiddleware(textSender services.TextSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("textSender", textSender)
		c.Next()
	}
}
