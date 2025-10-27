package cmd

import (
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/api"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/repository"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServeHTTP() {
	var dependency = dependencyIjection()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// routes := setupRoutes(e)

	api := e.Group("/api")

	// Healthcheck
	api.GET("/", dependency.HealthcheckAPI.HealthCheck)

	// Auth
	auth := api.Group("/v1/auth")
	auth.POST("/register", dependency.UserAPI.RegisterUser)

	e.Start(":" + helpers.GetEnv("PORT", "9000"))
}

func setupRoutes(e *echo.Echo) {
	// api := e.Group("/api")

	// // Healthcheck
	// api.GET("/", dependency.HealthcheckAPI.HealthCheck)

	// // Auth
	// auth := api.Group("/v1/auth")
	// auth.POST("/register", dependency.UserAPI.RegisterUser)
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
