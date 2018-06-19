package server

import (
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/services"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type API struct {
	Router      *gin.Engine
	Config      *viper.Viper
	Database    *mgo.Database
	EmailSender services.EmailSender
	Redis       *services.Redis
}
