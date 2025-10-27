package services

import (
	"context"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo interfaces.IUserRepository
}

func (s *UserService) Register(ctx context.Context, req *models.User, role string) (*models.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.Password = string(hashPassword)

	req.Role = role
	err = s.UserRepo.InsertNewUser(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := req
	resp.Password = ""
	return resp, nil
}
