package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models/dto"
)

type AuthService struct {
	authRepo interfaces.IAuthRepository
}

func NewAuthService(authRepo interfaces.IAuthRepository) interfaces.IAuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// Register handles user registration
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	existingUser, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to check email")
	}
	if existingUser != nil {
		return nil, helpers.ErrConflict("Email already registered")
	}

	// Check if username already exists
	existingUser, err = s.authRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to check username")
	}
	if existingUser != nil {
		return nil, helpers.ErrConflict("Username already taken")
	}

	// Check if phone already exists
	existPhone, err := s.authRepo.FindByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to check phone")
	}

	if existPhone != nil {
		return nil, helpers.ErrConflict("Phone already exist")
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to hash password")
	}

	// Parse DOB if provided
	var dob *time.Time
	if req.Dob != "" {
		parsedDob, err := time.Parse("2006-01-02", req.Dob)
		if err != nil {
			return nil, helpers.ErrBadRequest("Invalid date format for DOB. Use YYYY-MM-DD")
		}
		dob = &parsedDob
	}

	// Create user
	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		PhoneNumber:   req.PhoneNumber,
		FullName:      req.FullName,
		Address:       req.Address,
		Dob:           dob,
		Password:      hashedPassword,
		Role:          "user",
		EmailVerified: false,
		IsActive:      true,
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, helpers.ErrInternalServer("Failed to create user")
	}

	// Generate tokens
	accessToken, accessExpiry, err := helpers.GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to generate access token")
	}

	refreshToken, refreshExpiry, err := helpers.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to generate refresh token")
	}

	// Save session
	session := &models.UserSession{
		UserID:              user.ID,
		Token:               accessToken,
		RefreshToken:        refreshToken,
		TokenExpired:        accessExpiry,
		RefreshTokenExpired: refreshExpiry,
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, helpers.ErrInternalServer("Failed to create session")
	}

	// Prepare response
	response := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			FullName:    user.FullName,
			Address:     user.Address,
			Dob:         user.Dob,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
	}

	return response, nil
}

// Login handles user login
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email or username
	user, err := s.authRepo.FindByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to find user")
	}
	if user == nil {
		return nil, helpers.ErrUnauthorized("Invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, helpers.ErrUnauthorized("Account is deactivated")
	}

	// Verify password
	if err := helpers.ComparePassword(user.Password, req.Password); err != nil {
		return nil, helpers.ErrUnauthorized("Invalid credentials")
	}

	// Generate tokens
	accessToken, accessExpiry, err := helpers.GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to generate access token")
	}

	refreshToken, refreshExpiry, err := helpers.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to generate refresh token")
	}

	// Delete old sessions and create new one
	_ = s.authRepo.DeleteSessionsByUserID(ctx, user.ID)

	session := &models.UserSession{
		UserID:              user.ID,
		Token:               accessToken,
		RefreshToken:        refreshToken,
		TokenExpired:        accessExpiry,
		RefreshTokenExpired: refreshExpiry,
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, helpers.ErrInternalServer("Failed to create session")
	}

	// Prepare response
	response := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			FullName:    user.FullName,
			Address:     user.Address,
			Dob:         user.Dob,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
	}

	return response, nil
}

// RefreshToken handles token refresh
func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := helpers.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, helpers.ErrUnauthorized("Invalid refresh token")
	}

	// Find session
	session, err := s.authRepo.FindSessionByUserID(ctx, claims.UserID)

	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to find session")
	}
	if session == nil || session.RefreshToken != req.RefreshToken {
		return nil, helpers.ErrUnauthorized("Invalid refresh token")
	}

	// Check if refresh token is expired
	if time.Now().After(session.RefreshTokenExpired) {
		return nil, helpers.ErrUnauthorized("Refresh token expired")
	}

	// Get user details
	user, err := s.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to find user")
	}
	if user == nil {
		return nil, helpers.ErrUnauthorized("User not found")
	}

	if !user.IsActive {
		return nil, helpers.ErrUnauthorized("Account is deactivated")
	}

	// Generate new tokens
	accessToken, accessExpiry, err := helpers.GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to generate access token")
	}

	newRefreshToken, refreshExpiry, err := helpers.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to generate refresh token")
	}

	// Update session
	session.Token = accessToken
	session.RefreshToken = newRefreshToken
	session.TokenExpired = accessExpiry
	session.RefreshTokenExpired = refreshExpiry

	if err := s.authRepo.UpdateSession(ctx, session); err != nil {
		return nil, helpers.ErrInternalServer("Failed to update session")
	}

	// Prepare response
	response := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			FullName:    user.FullName,
			Address:     user.Address,
			Dob:         user.Dob,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiry,
	}

	return response, nil
}

// ForgotPassword handles password reset request
func (s *AuthService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	// Find user by email
	user, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return helpers.ErrInternalServer("Failed to find user")
	}
	if user == nil {
		// Don't reveal if email exists or not
		return nil
	}

	// Generate reset token
	resetToken, err := helpers.GenerateRandomToken(32)
	if err != nil {
		return helpers.ErrInternalServer("Failed to generate reset token")
	}

	// Set expiry time (1 hour)
	expiry := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	// Save reset token
	if err := s.authRepo.SaveResetToken(ctx, user.ID, resetToken, expiry); err != nil {
		return helpers.ErrInternalServer("Failed to save reset token")
	}

	// TODO: Send email with reset token
	// For now, we'll just log it (in production, use email service)
	fmt.Printf("Password reset token for %s: %s\n", user.Email, resetToken)

	return nil
}

// ResetPassword handles password reset
func (s *AuthService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// Find user by reset token
	user, err := s.authRepo.FindByResetToken(ctx, req.Token)
	if err != nil {
		return helpers.ErrInternalServer("Failed to verify reset token")
	}
	if user == nil {
		return helpers.ErrBadRequest("Invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return helpers.ErrInternalServer("Failed to hash password")
	}

	// Update password
	if err := s.authRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return helpers.ErrInternalServer("Failed to update password")
	}

	// Clear reset token
	if err := s.authRepo.ClearResetToken(ctx, user.ID); err != nil {
		return helpers.ErrInternalServer("Failed to clear reset token")
	}

	// Delete all sessions (force re-login)
	_ = s.authRepo.DeleteSessionsByUserID(ctx, user.ID)

	return nil
}

// ChangePassword handles password change for authenticated user
func (s *AuthService) ChangePassword(ctx context.Context, userID int, req *dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.authRepo.FindByID(ctx, userID)
	if err != nil {
		return helpers.ErrInternalServer("Failed to find user")
	}
	if user == nil {
		return helpers.ErrNotFound("User not found")
	}

	// Verify old password
	if err := helpers.ComparePassword(user.Password, req.OldPassword); err != nil {
		return helpers.ErrBadRequest("Invalid old password")
	}

	// Hash new password
	hashedPassword, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return helpers.ErrInternalServer("Failed to hash password")
	}

	// Update password
	if err := s.authRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return helpers.ErrInternalServer("Failed to update password")
	}

	// Delete all sessions except current (force re-login on other devices)
	_ = s.authRepo.DeleteSessionsByUserID(ctx, userID)

	return nil
}

// Logout handles user logout
func (s *AuthService) Logout(ctx context.Context, userID int, token string) error {
	// Delete session by token
	if err := s.authRepo.DeleteSession(ctx, token); err != nil {
		return helpers.ErrInternalServer("Failed to logout")
	}

	return nil
}

// GetProfile retrieves user profile
func (s *AuthService) GetProfile(ctx context.Context, userID int) (*dto.UserResponse, error) {
	user, err := s.authRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, helpers.ErrInternalServer("Failed to find user")
	}
	if user == nil {
		return nil, helpers.ErrNotFound("User not found")
	}

	response := &dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
		Address:     user.Address,
		Dob:         user.Dob,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}

	return response, nil
}
