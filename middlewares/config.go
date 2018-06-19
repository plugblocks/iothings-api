package middlewares

import (
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/config"
	"gopkg.in/gin-gonic/gin.v1"
)

func ConfigMiddleware(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ToContext(c, config.New(viper))
		c.Next()
	}
}
