package server

import (
	"gitlab.com/plugblocks/iothings-api/config"
	"net/http"
	"time"

	"gitlab.com/plugblocks/iothings-api/controllers"
	"gitlab.com/plugblocks/iothings-api/middlewares"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the " + config.GetString(c, "mail_sender_name") + " API."})
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
	router.Use(middlewares.TextMiddleware(a.TextSender))

	if a.Config.GetBool("plan_check") == true {
		router.Use(middlewares.PlanMiddleware())
	}

	//customerAuthMiddleware := middlewares.CustomerAuthMiddleware()
	authMiddleware := middlewares.AuthMiddleware() //User
	adminMiddleware := middlewares.AdminMiddleware()

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			authentication.POST("/", authController.UserAuthentication)
		}

		users := v1.Group("/users")
		{
			userController := controllers.NewUserController()
			//v1.POST("/reset_password", userController.ResetPasswordRequest)
			users.GET("/:id/activate/:activationKey", userController.ActivateUser)
			users.Use(authMiddleware)
			users.GET("/:id/organization", userController.GetUserOrganization)
			users.Use(adminMiddleware)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.DELETE("/:id", userController.DeleteUser)
			users.GET("/", userController.GetUsers)
			users.PUT("/:id/assign/:organization_id", userController.AssignOrganization)
		}

		customers := v1.Group("/customers")
		{
			customerController := controllers.NewCustomerController()
			customers.Use(authMiddleware)
			customers.POST("/", customerController.CreateCustomer)
			customers.GET("/:id", customerController.GetCustomerById)
			customers.DELETE("/:id", customerController.DeleteCustomer)
			customers.PUT("/:id", customerController.EditCustomer)
			customers.GET("/", customerController.GetCustomers)
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

		fleets := v1.Group("/fleets")
		{
			fleetsController := controllers.NewFleetController()
			//TODO: DANGER: Protect by auth device GeoJSON
			fleets.GET("/:id/locations/geojson", fleetsController.GetFleetGeoJSON)
			fleets.Use(authMiddleware)
			fleets.GET("/", fleetsController.GetFleets)
			fleets.POST("/", fleetsController.CreateFleet)
			fleets.PUT("/:id", fleetsController.UpdateFleet)
			fleets.GET("/:id", fleetsController.GetFleetById)
			fleets.GET("/:id/devices", fleetsController.GetDevicesFromFleet)
			fleets.POST("/:id/:deviceId", fleetsController.AddDeviceToFleet)
			fleets.DELETE("/:id", fleetsController.DeleteFleet)
		}

		devices := v1.Group("/devices")
		{
			deviceController := controllers.NewDeviceController()
			//TODO: DANGER: Protect by auth device GeoJSON
			devices.GET("/:id/locations/geojson", deviceController.GetDeviceGeoJSON)
			devices.Use(authMiddleware)
			devices.GET("/:id/location/:source", deviceController.GetDeviceLastLocation)
			devices.GET("/:id/locations", deviceController.GetDeviceGeolocations)
			devices.GET("/:id/messages", deviceController.GetDeviceMessages)
			devices.GET("/", deviceController.GetDevices)
			devices.POST("/", deviceController.CreateDevice)
			devices.PUT("/:id", deviceController.UpdateDevice)
			devices.GET("/:id", deviceController.GetDevice)
			devices.DELETE("/:id", deviceController.DeleteDevice)
		}

		filters := v1.Group("/filters")
		{
			deviceController := controllers.NewDeviceController()
			filters.Use(authMiddleware)
			filters.GET("/devices/available", deviceController.GetAvailableDevices)
		}

		alerts := v1.Group("/alerts")
		{
			alertController := controllers.NewAlertController()
			alerts.Use(authMiddleware)
			alerts.GET("/", alertController.GetAlerts)
			alerts.GET("/:id", alertController.GetAlert)
			alerts.POST("/", alertController.CreateAlert)
			alerts.PUT("/:id", alertController.UpdateAlert)
			alerts.DELETE("/:id", alertController.DeleteAlert)
			alerts.POST("/fleet/:fleetId", alertController.GetFleetAlerts)
			alerts.POST("/device/:deviceId", alertController.GetDeviceAlerts)
			alerts.POST("/sms/:userId", alertController.TextTest)
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
			organizations.GET("/:id/subscription", organizationsController.GetOrganizationSubscription)
			organizations.DELETE("/:id", organizationsController.DeleteOrganization)
			organizations.GET("/:id/users", organizationsController.GetUsers)
		}

		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.Use(authMiddleware)
			subscriptions.Use(adminMiddleware)
			subscriptionsController := controllers.NewSubscriptionController()
			subscriptions.GET("/", subscriptionsController.GetSubscriptions)
			subscriptions.POST("/", subscriptionsController.CreateSubscription)
			subscriptions.PUT("/:id", subscriptionsController.UpdateSubscription)
			subscriptions.GET("/:id", subscriptionsController.GetSubscription)
			subscriptions.DELETE("/:id", subscriptionsController.DeleteSubscription)
		}

		orders := v1.Group("/orders")
		{
			orders.Use(authMiddleware)
			ordersController := controllers.NewOrderController()
			orders.GET("/", ordersController.GetOrders)
			orders.POST("/", ordersController.CreateOrder)
			orders.PUT("/:id", ordersController.EditOrder)
			orders.GET("/:id", ordersController.GetOrderById)
			orders.DELETE("/:id", ordersController.DeleteOrder)
			orders.PUT("/:id/terminate", ordersController.TerminateOrder)
			orders.GET("/:id/locations", ordersController.GetOrderGeolocations)
			orders.GET("/:id/matching_map", ordersController.GetMatchingMap)
		}

		warehouses := v1.Group("/warehouses")
		{
			warehouses.Use(authMiddleware)
			warehousesController := controllers.NewWarehouseController()
			warehouses.GET("/", warehousesController.GetWarehouses)
			warehouses.POST("/", warehousesController.CreateWarehouse)
			warehouses.PUT("/:id", warehousesController.EditWarehouse)
			warehouses.GET("/:id", warehousesController.GetWarehouseById)
			warehouses.DELETE("/:id", warehousesController.DeleteWarehouse)
		}

		observations := v1.Group("/observations")
		{
			observationController := controllers.NewObservationController()
			observations.POST("/new", observationController.CreateObservation)
			observations.Use(authMiddleware)
			observations.GET("/device/:deviceId", observationController.GetDeviceObservations)
			observations.GET("/fleet/:fleetId", observationController.GetFleetObservations)
			observations.DELETE("/:id", observationController.DeleteObservation)
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
			geolocations.GET("/user/fleets", fleetsController.GetUserFleetsGeoJSON)
		}
	}
}
