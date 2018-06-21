package middlewares

import (
	"gitlab.com/plugblocks/iothings-api/services"

	"github.com/gin-gonic/gin"
)

func RedisMiddleware(redis *services.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redis)
		c.Next()
	}
}
