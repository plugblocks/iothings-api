package server

import (
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

type API struct {
	Router      *gin.Engine
	Config      *viper.Viper
	Database    *mgo.Database
	EmailSender services.EmailSender
	Redis       *services.Redis
}
