package middlewares

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/plugblocks/iothings-api/helpers"
	"gitlab.com/plugblocks/iothings-api/store"
	"gopkg.in/gin-gonic/gin.v1"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := store.Current(c)

		if !user.Admin {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("admin_required", "The user is not administrator", errors.New("The user is not administrator")))
			return
		}

		c.Next()
	}
}
