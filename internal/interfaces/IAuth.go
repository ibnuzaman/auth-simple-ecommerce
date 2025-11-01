package interfaces

import (
	"context"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models/dto"
)

type IAuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID int, req *dto.ChangePasswordRequest) error
	Logout(ctx context.Context, userID int, token string) error
	GetProfile(ctx context.Context, userID int) (*dto.UserResponse, error)
}

type IAuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmailOrUsername(ctx context.Context, emailOrUsername string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	SaveResetToken(ctx context.Context, userID int, token string, expiry string) error
	FindByResetToken(ctx context.Context, token string) (*models.User, error)
	ClearResetToken(ctx context.Context, userID int) error

	// Session management
	CreateSession(ctx context.Context, session *models.UserSession) error
	FindSessionByToken(ctx context.Context, token string) (*models.UserSession, error)
	FindSessionByUserID(ctx context.Context, userID int) (*models.UserSession, error)
	UpdateSession(ctx context.Context, session *models.UserSession) error
	DeleteSession(ctx context.Context, token string) error
	DeleteSessionsByUserID(ctx context.Context, userID int) error
}
