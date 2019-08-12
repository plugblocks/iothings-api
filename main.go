package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"gitlab.com/plugblocks/iothings-api/server"
	"gitlab.com/plugblocks/iothings-api/services"
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
	api.SetupSeeds()

	if api.Config.GetBool("plan_check") == true {
		fmt.Println("Plan check enabled")
		services.CheckSubscription(api.Router, api.Config, &gin.Context{})
		cron := cron.New()
		cron.AddFunc("@every 1m", func() {
			services.CheckSubscription(api.Router, api.Config, &gin.Context{})
		})
		cron.Start()
	}

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
