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
			users.Use(adminMiddleware)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.GET("/", userController.GetUsers)
		}

		fleets := v1.Group("/fleets")
		{
			fleets.Use(authMiddleware)
			fleetsController := controllers.NewFleetController()
			fleets.GET("/", fleetsController.GetFleets)
			fleets.POST("/", fleetsController.CreateFleet)
			fleets.PUT("/:id", fleetsController.EditFleet)
			fleets.GET("/:id", fleetsController.GetFleetById)
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
			devices.Use(authMiddleware)
			deviceController := controllers.NewDeviceController()
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
			organizations.GET("/:id/users", organizationsController.GetOrganizationUsers)
			organizations.DELETE("/:id", organizationsController.DeleteOrganization)
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
			observations.GET("/device/:id/:type", observationController.GetDeviceObservations)
			observations.GET("/device/:id/:type/latest", observationController.GetDeviceLatestObservation)
			observations.GET("/fleet/:id/:type", observationController.GetFleetObservations)
			observations.GET("/fleet/:id/:type/latest", observationController.GetFleetLatestObservation)
			observations.GET("/fleets/:type", observationController.GetAllFleetsObservations)
			observations.GET("/fleets/:type/latest", observationController.GetAllFleetsLatestObservation)
		}

		sigfox := v1.Group("/sigfox")
		{
			sigfoxController := controllers.NewSigfoxController()
			sigfox.POST("/message", sigfoxController.CreateSigfoxMessage)
			sigfox.POST("/location", sigfoxController.CreateSigfoxLocation)
		}
	}
}
