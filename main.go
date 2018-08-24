package main

import (
	"gitlab.com/plugblocks/iothings-api/server"
	"gitlab.com/plugblocks/iothings-api/services"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
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
	api.EmailSender = services.NewEmailSender(api.Config)

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
	api.SetupSeeds()

	services.CheckSubscription(api.Config)

	cron := cron.New()
	cron.AddFunc("@every 6h", func() {
		services.CheckSubscription(api.Config)
	})
	cron.Start()

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
