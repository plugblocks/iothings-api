package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/config"
)

func ConfigMiddleware(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ToContext(c, config.New(viper))
		c.Next()
	}
}
