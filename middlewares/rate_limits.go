package middlewares

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/services"
	"github.com/gin-gonic/gin"
)

func RateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := services.GetRedis(c).Pool.Get()
		defer conn.Close()

		if !config.GetBool(c, "rate_limit_activated") {
			return
		}

		unixTime := int64(time.Now().Unix())
		keyName := c.ClientIP() + ":" + strconv.FormatInt(unixTime, 10)

		var count int = -1
		data, err := conn.Do("GET", keyName)
		if err != nil {
			return
		}

		if data != nil {
			if err := json.Unmarshal(data.([]byte), &count); err != nil {
				return
			}
		}

		if count != -1 && count >= config.GetInt(c, "rate_limit_requests_per_second") {
			c.AbortWithError(http.StatusTooManyRequests, helpers.ErrorWithCode("too_many_requests", "You sent too many requests over the last second.", err))
		} else {
			conn.Do("INCR", keyName)
			conn.Do("EXPIRE", keyName, 10)
		}

		c.Next()
	}
}
