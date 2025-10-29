package cmd

import (
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/api"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/repository"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	_ "github.com/ibnuzaman/auth-simple-ecommerce.git/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ServeHTTP() {
	var dependency = dependencyIjection()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	api := e.Group("/api")
	api.GET("/swagger/*", echoSwagger.WrapHandler)

	// Healthcheck
	api.GET("/", dependency.HealthcheckAPI.HealthCheck)

	// Auth
	auth := api.Group("/v1/auth")
	auth.POST("/register", dependency.UserAPI.RegisterUser)

	if err := e.Start(":" + helpers.GetEnv("PORT", "9000")); err != nil {
		logrus.Info("Failed to connect app", err)
		return
	}

}

type Dependency struct {
	HealthcheckAPI *api.HealthCheckAPI

	UserAPI interfaces.IUserAPI
}

func dependencyIjection() Dependency {
	userRepo := &repository.UserRepository{
		DB: helpers.DB,
	}

	userService := &services.UserService{
		UserRepo: userRepo,
	}

	userAPI := &api.UserAPI{
		UserService: *userService,
	}

	return Dependency{
		HealthcheckAPI: &api.HealthCheckAPI{},
		UserAPI:        userAPI,
	}
}
