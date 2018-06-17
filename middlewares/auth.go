package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"gitlab.com/plugblocks/iothings-api/models"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/store"
	"gopkg.in/gin-gonic/gin.v1"
	"gitlab.com/plugblocks/iothings-api/config"
	"gitlab.com/plugblocks/iothings-api/helpers"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenReader := c.Request.Header.Get("Authorization")

		authHeaderParts := strings.Split(tokenReader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			c.AbortWithError(http.StatusBadRequest, errors.New("Authorization header format must be Bearer {token}"))
			return
		}

		encodedKey := []byte(config.GetString(c, "rsa_private"))
		claims, err := helpers.ValidateJwtToken(authHeaderParts[1], encodedKey, "access")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_token", "the given token is invalid", err))
			return
		}

		user := &models.User{}

		// Gets the user from the redis store
		err = services.GetRedis(c).GetValueForKey(claims["id"].(string), &user)
		if err != nil {
			user, _ = store.FindUserById(c, claims["id"].(string))
			services.GetRedis(c).SetValueForKey(user.Id, &user)
		}

		c.Set(store.CurrentKey, user)

		c.Next()
	}
}
