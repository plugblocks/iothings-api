package server

import (
	"os"

	"github.com/joho/godotenv"
)

func (a *API) SetupViper() error {

	filename := ".env"
	switch os.Getenv("IOTHINGS_ENV") {
	case "testing":
		filename = "../.env.testing"
	case "prod":
		filename = ".env.prod"
	}

	godotenv.Overload(filename)

	a.Config.SetEnvPrefix("iothings")
	a.Config.AutomaticEnv()

	a.SetupViperDefaults()

	return nil
}

func (a *API) SetupViperDefaults() {
	a.Config.SetDefault("rate_limit_requests_per_second", 5)
	a.Config.SetDefault("rate_limit_activated", true)
}
