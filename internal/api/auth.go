package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/interfaces"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models/dto"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService interfaces.IAuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService interfaces.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} helpers.BaseResponse{data=dto.AuthResponse}
// @Failure 400 {object} helpers.BaseResponse
// @Failure 409 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusCreated, "Registration successful", response)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} helpers.BaseResponse{data=dto.AuthResponse}
// @Failure 400 {object} helpers.BaseResponse
// @Failure 401 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Login successful", response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} helpers.BaseResponse{data=dto.AuthResponse}
// @Failure 400 {object} helpers.BaseResponse
// @Failure 401 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	response, err := h.authService.RefreshToken(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Token refreshed successfully", response)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset token to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "User email"
// @Success 200 {object} helpers.BaseResponse{data=dto.MessageResponse}
// @Failure 400 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req dto.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	err := h.authService.ForgotPassword(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Password reset instructions sent to your email", dto.MessageResponse{
		Message: "If the email exists, you will receive password reset instructions",
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} helpers.BaseResponse{data=dto.MessageResponse}
// @Failure 400 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	err := h.authService.ResetPassword(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Password reset successful", dto.MessageResponse{
		Message: "Your password has been reset successfully. Please login with your new password",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} helpers.BaseResponse{data=dto.MessageResponse}
// @Failure 400 {object} helpers.BaseResponse
// @Failure 401 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return helpers.ResponseHttp(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return helpers.ResponseHttp(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	err := h.authService.ChangePassword(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Password changed successfully", dto.MessageResponse{
		Message: "Your password has been changed successfully",
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.BaseResponse{data=dto.MessageResponse}
// @Failure 401 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return helpers.ResponseHttp(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	token, ok := c.Get("token").(string)
	if !ok {
		return helpers.ResponseHttp(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	err := h.authService.Logout(c.Request().Context(), userID, token)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Logout successful", dto.MessageResponse{
		Message: "You have been logged out successfully",
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.BaseResponse{data=dto.UserResponse}
// @Failure 401 {object} helpers.BaseResponse
// @Failure 500 {object} helpers.BaseResponse
// @Router /v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return helpers.ResponseHttp(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	response, err := h.authService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return helpers.ResponseHttp(c, http.StatusOK, "Profile retrieved successfully", response)
}
