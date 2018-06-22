package middlewares

import (
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/config"
	"github.com/gin-gonic/gin"
)

func ConfigMiddleware(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ToContext(c, config.New(viper))
		c.Next()
	}
}
