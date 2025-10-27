package api

import (
	"net/http"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/labstack/echo/v4"
)

type home struct {
	Title string `json:"title"`
}

type HealthCheckAPI struct {
}

func (api *HealthCheckAPI) HealthCheck(e echo.Context) error {
	return helpers.ResponseHttp(e, http.StatusOK, "Healthty", home{Title: "Welcome to the auth service simple ecommerce"})
}
