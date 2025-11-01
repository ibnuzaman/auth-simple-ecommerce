package api

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/constants"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models/dto"
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
var validate = validator.New()

func (api *UserAPI) RegisterUser(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, constants.ErrBadRequest, nil)
	}
	if err := validate.Struct(req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, constants.ErrBadRequest, echo.Map{
			"details": err.Error(),
		})
	}

	resp, err := api.UserService.Register(c.Request().Context(), req, constants.RoleCustomer)
	if err != nil {
		var cf *constants.ConflictError
		if errors.As(err, &cf) {
			return helpers.ResponseHttp(c, http.StatusConflict, constants.ErrDuplicate, echo.Map{
				"field":   cf.Field,
				"message": cf.Field + " already in use",
			})
		}
		return helpers.ResponseHttp(c, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	// 201 Created lebih tepat untuk registrasi
	return helpers.ResponseHttp(c, http.StatusCreated, constants.SuccessMessage, resp)
}
