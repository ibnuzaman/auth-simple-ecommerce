package cmd

import (
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServeHTTP() {
	d := dependencyIjection()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.GET("/health", d.HealthcheckAPI.HealthCheck)

	e.Start(":" + helpers.GetEnv("PORT", "9000"))
}

type Dependency struct {
	HealthcheckAPI *api.HealthCheckAPI
}

func dependencyIjection() Dependency {
	return Dependency{
		HealthcheckAPI: &api.HealthCheckAPI{},
	}
}
