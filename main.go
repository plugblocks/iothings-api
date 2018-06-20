package main

import (
	"gitlab.com/plugblocks/iothings-api/server"
	"gitlab.com/plugblocks/iothings-api/services"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	fmt.Println("Config: Host Address" + api.Config.GetString("host_address"))
	fmt.Println("Config: Databases: Mongo:" + api.Config.GetString("DB_HOST") + "\t Redis:" + api.Config.GetString("REDIS_HOST"))

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

	// Seeds setup
	api.SetupSeeds()

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
