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

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new customer user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User				true	"User registration details"
//	@Success		200		{object}	helpers.BaseResponse	"Successfully registered user"
//	@Failure		400		{object}	helpers.BaseResponse	"Bad request - invalid input"
//	@Failure		500		{object}	helpers.BaseResponse	"Internal server error"
//	@Router			/api/v1/auth/register [post]
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

	resp, err := api.UserService.Register(e.Request().Context(), &req, constants.RoleCustomer)
	if err != nil {
		log.Error("failed to register: ", err)
		return helpers.ResponseHttp(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.ResponseHttp(e, http.StatusOK, constants.SuccessMessage, resp)
}
