package server

import (
	"net/http"
	"time"

	"gitlab.com/plugblocks/iothings-api/controllers"
	"gitlab.com/plugblocks/iothings-api/middlewares"

	"gopkg.in/gin-gonic/gin.v1"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the base API."})
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
	router.Use(middlewares.RedisMiddleware(a.Redis))

	router.Use(middlewares.EmailMiddleware(a.EmailSender))
	router.Use(middlewares.RateMiddleware())

	authMiddleware := middlewares.AuthMiddleware()
	adminMiddleware := middlewares.AdminMiddleware()

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		userController := controllers.NewUserController()
		//v1.POST("/reset_password", userController.ResetPasswordRequest)
		users := v1.Group("/users")
		{
			users.POST("/", userController.CreateUser)
			users.GET("/:id/activate/:activationKey", userController.ActivateUser)
			users.Use(authMiddleware)
			users.GET("/:id", userController.GetUser)
			users.Use(adminMiddleware)
			users.GET("/", userController.GetUsers)
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

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			authentication.POST("/", authController.Authentication)
		}
	}
}
