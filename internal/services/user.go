package services

import (
	"context"
	"errors"
	"time"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/constants"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo interfaces.IUserRepository
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest, role string) (*dto.RegisterResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var dobPtr *time.Time
	if req.Dob != "" {
		t, err := time.Parse("2006-01-02", req.Dob)
		if err != nil {
			return nil, constants.ErrFailedBadRequest
		}
		dobPtr = &t
	}

	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FullName:    req.FullName,
		Address:     req.Address,
		Dob:         dobPtr,
		Password:    string(hashPassword),
		Role:        role,
	}

	if err := s.UserRepo.InsertNewUser(ctx, user); err != nil {
		if errors.Is(err, constants.ErrConflict) {
			return nil, err
		}
		return nil, err
	}

	resp := &dto.RegisterResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
		Address:     user.Address,
		Dob:         req.Dob,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}
	return resp, nil
}
