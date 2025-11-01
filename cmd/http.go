package cmd

import (
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/api"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	appMiddleware "github.com/ibnuzaman/auth-simple-ecommerce.git/internal/middleware"
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

	// Custom error handler
	e.HTTPErrorHandler = appMiddleware.ErrorHandler

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api")
	api.GET("/swagger/*", echoSwagger.WrapHandler)

	// Healthcheck
	api.GET("/health", dependency.HealthcheckAPI.HealthCheck)

	// Auth routes (public)
	auth := api.Group("/v1/auth")
	auth.POST("/register", dependency.AuthAPI.Register)
	auth.POST("/login", dependency.AuthAPI.Login)
	auth.POST("/refresh", dependency.AuthAPI.RefreshToken)
	auth.POST("/forgot-password", dependency.AuthAPI.ForgotPassword)
	auth.POST("/reset-password", dependency.AuthAPI.ResetPassword)

	// Auth routes (protected)
	authProtected := api.Group("/v1/auth")
	authProtected.Use(appMiddleware.JWTMiddleware())
	authProtected.POST("/logout", dependency.AuthAPI.Logout)
	authProtected.POST("/change-password", dependency.AuthAPI.ChangePassword)
	authProtected.GET("/profile", dependency.AuthAPI.GetProfile)

	// User routes (protected) - example
	users := api.Group("/v1/users")
	users.Use(appMiddleware.JWTMiddleware())
	// Add user routes here if needed

	if err := e.Start(":" + helpers.GetEnv("PORT", "9000")); err != nil {
		logrus.Info("Failed to connect app", err)
		return
	}

}

type Dependency struct {
	HealthcheckAPI *api.HealthCheckAPI
	AuthAPI        *api.AuthHandler
	UserAPI        interfaces.IUserAPI
}

func dependencyIjection() Dependency {
	// Auth dependencies
	authRepo := repository.NewAuthRepository(helpers.DB)
	authService := services.NewAuthService(authRepo)
	authAPI := api.NewAuthHandler(authService)

	// User dependencies
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
		AuthAPI:        authAPI,
		UserAPI:        userAPI,
	}
}
