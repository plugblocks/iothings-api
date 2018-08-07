package server

import (
	"net/http"
	"time"

	"gitlab.com/plugblocks/iothings-api/controllers"
	"gitlab.com/plugblocks/iothings-api/middlewares"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the iothings API."})
}

func (a *API) SetupRouter() {
	router := a.Router

	router.Use(middlewares.ErrorMiddleware())

	router.Use(middlewares.CorsMiddleware(middlewares.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.Use(middlewares.StoreMiddleware(a.Database))
	router.Use(middlewares.ConfigMiddleware(a.Config))
	router.Use(middlewares.EmailMiddleware(a.EmailSender))

	authMiddleware := middlewares.AuthMiddleware()
	adminMiddleware := middlewares.AdminMiddleware()

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		userController := controllers.NewUserController()
		//v1.POST("/reset_password", userController.ResetPasswordRequest)
		users := v1.Group("/users")
		{
			users.GET("/:id/activate/:activationKey", userController.ActivateUser)
			users.Use(authMiddleware)
			users.GET("/:id/organization", userController.GetUserOrganization)
			users.Use(adminMiddleware)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.GET("/", userController.GetUsers)
			users.PUT("/:id/assign/:organization_id", userController.AssignOrganization)
		}

		fleets := v1.Group("/fleets")
		{
			fleetsController := controllers.NewFleetController()
			//TODO: DANGER: Protect by auth device GeoJSON
			fleets.GET("/:id/locations/geojson", fleetsController.GetFleetGeoJSON)
			fleets.Use(authMiddleware)
			fleets.GET("/", fleetsController.GetFleets)
			fleets.POST("/", fleetsController.CreateFleet)
			fleets.PUT("/:id", fleetsController.EditFleet)
			fleets.GET("/:id", fleetsController.GetFleetById)
			fleets.POST("/:id/:deviceId", fleetsController.AddDeviceToFleet)
			fleets.DELETE("/:id", fleetsController.DeleteFleet)
		}

		groups := v1.Group("/groups")
		{
			groups.Use(authMiddleware)
			groupsController := controllers.NewGroupController()
			groups.GET("/", groupsController.GetGroups)
			groups.POST("/", groupsController.CreateGroup)
			groups.PUT("/:id", groupsController.EditGroup)
			groups.GET("/:id", groupsController.GetGroupById)
			groups.DELETE("/:id", groupsController.DeleteGroup)
		}

		devices := v1.Group("/devices")
		{
			deviceController := controllers.NewDeviceController()
			//TODO: DANGER: Protect by auth device GeoJSON
			devices.GET("/:id/locations/geojson", deviceController.GetDeviceGeoJSON)
			devices.Use(authMiddleware)
			devices.GET("/", deviceController.GetDevices)
			devices.POST("/", deviceController.CreateDevice)
			devices.PUT("/:id", deviceController.UpdateDevice)
			devices.GET("/:id", deviceController.GetDevice)
			devices.DELETE("/:id", deviceController.DeleteDevice)
		}

		organizations := v1.Group("/organizations")
		{
			organizations.Use(authMiddleware)
			organizations.Use(adminMiddleware)
			organizationsController := controllers.NewOrganizationController()
			organizations.GET("/", organizationsController.GetOrganizations)
			organizations.POST("/", organizationsController.CreateOrganization)
			organizations.PUT("/:id", organizationsController.UpdateOrganization)
			organizations.GET("/:id", organizationsController.GetOrganizationById)
			organizations.DELETE("/:id", organizationsController.DeleteOrganization)
			organizations.GET("/:id/users", organizationsController.GetUsers)
		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			authentication.POST("/", authController.Authentication)
		}

		observations := v1.Group("/observations")
		{
			observationController := controllers.NewObservationController()
			observations.POST("/new", observationController.CreateObservation)
			observations.Use(authMiddleware)
			observations.GET("/device/:deviceId", observationController.GetDeviceObservations)
			observations.GET("/fleet/:fleetId", observationController.GetFleetObservations)
		}

		sigfox := v1.Group("/sigfox")
		{
			sigfoxController := controllers.NewSigfoxController()
			sigfox.POST("/message", sigfoxController.CreateSigfoxMessage)
			sigfox.POST("/location", sigfoxController.CreateSigfoxLocation)
			//sigfox.GET("/locations", sigfoxController.GetSigfoxLocations)
			//sigfox.GET("/locations/geojson", sigfoxController.GetGeoJSON)
		}

		geolocations := v1.Group("/geolocations")
		{
			geolocationController := controllers.NewGeolocationController()
			fleetsController := controllers.NewFleetController()
			geolocations.GET("/fleets", fleetsController.GetFleetsGeoJSON)
			geolocations.POST("/", geolocationController.CreateGeolocation)
			geolocations.DELETE("/:id", geolocationController.DeleteGeolocation)

			geolocations.Use(authMiddleware)
			geolocations.Use(adminMiddleware)
			geolocations.GET("/user/fleets", fleetsController.GetUserFleetsGeoJSON)
		}
	}
}
