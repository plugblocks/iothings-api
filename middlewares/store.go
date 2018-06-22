package middlewares

import (
	"github.com/gin-gonic/gin"
	mgo "github.com/globalsign/mgo"
	"gitlab.com/plugblocks/iothings-api/store"
	"gitlab.com/plugblocks/iothings-api/store/mongodb"
)

func StoreMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, mongodb.New(db))
		c.Next()
	}
}
