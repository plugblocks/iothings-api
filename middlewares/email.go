package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/plugblocks/iothings-api/services"
)

func EmailMiddleware(emailSender services.EmailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("emailSender", emailSender)
		c.Next()
	}
}
