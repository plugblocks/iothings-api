package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/server"
	"gitlab.com/plugblocks/iothings-api/services"
	"gitlab.com/plugblocks/iothings-api/utils"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	// Configuration setup
	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	// Email sender setup
	api.EmailSender = services.NewEmailSender(api.Config)
	api.TextSender = services.NewTextSender(api.Config)

	// Database setup
	session, err := api.SetupDatabase()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}

	// Seeds setup
	utils.CheckErr(api.SetupSeeds())

	// Router setup
	api.SetupRouter()
	utils.CheckErr(api.Router.Run(api.Config.GetString("host_address")))
}
