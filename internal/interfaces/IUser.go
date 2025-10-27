package interfaces

import (
	"context"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"github.com/labstack/echo/v4"
)

type IUserRepository interface {
	InsertNewUser(ctx context.Context, user *models.User) error
}

type IUserService interface {
	Register(ctx context.Context, req *models.User, role string) (*models.User, error)
}

type IUserAPI interface {
	RegisterUser(e echo.Context) error
}
