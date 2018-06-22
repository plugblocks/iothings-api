package main

import (
	"gitlab.com/plugblocks/iothings-api/server"
	"gitlab.com/plugblocks/iothings-api/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	// Configuration setup
	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	// Email sender setup
	api.EmailSender = services.NewSendGridEmailSender(api.Config)

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

	// Stripe setup
	//services.SetStripeKeyAndBackend(api.Config)

	// Redis setup
	api.SetupRedis()

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
