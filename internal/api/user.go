package api

import (
	"net/http"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/constants"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/services"
	"github.com/labstack/echo/v4"
)

type UserAPI struct {
	UserService services.UserService
}

func (api *UserAPI) RegisterUser(e echo.Context) error {
	var (
		log = helpers.Logger
	)
	req := models.User{}

	if err := e.Bind(&req); err != nil {
		log.Error("failed to parse request: ", err)
		return helpers.ResponseHttp(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request: ", err)
		return helpers.ResponseHttp(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	resp, err := api.UserService.Register(e.Request().Context(), &req, "customer")
	if err != nil {
		log.Error("failed to register: ", err)
		return helpers.ResponseHttp(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.ResponseHttp(e, http.StatusOK, constants.SuccessMessage, resp)
}
