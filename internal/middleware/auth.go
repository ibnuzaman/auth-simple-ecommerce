package middleware

import (
	"net/http"
	"strings"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates JWT token and adds user info to context
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return helpers.ResponseHttp(c, http.StatusUnauthorized, "Missing authorization header", nil)
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return helpers.ResponseHttp(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			}

			token := parts[1]

			// Validate token
			claims, err := helpers.ValidateToken(token)
			if err != nil {
				return helpers.ResponseHttp(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			}

			// Add user info to context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Set("token", token)

			return next(c)
		}
	}
}

// OptionalJWTMiddleware validates JWT token if present but doesn't require it
func OptionalJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				// Check if it's a Bearer token
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token := parts[1]

					// Validate token
					claims, err := helpers.ValidateToken(token)
					if err == nil {
						// Add user info to context
						c.Set("user_id", claims.UserID)
						c.Set("email", claims.Email)
						c.Set("username", claims.Username)
						c.Set("role", claims.Role)
						c.Set("token", token)
					}
				}
			}

			return next(c)
		}
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("role").(string)
			if !ok {
				return helpers.ResponseHttp(c, http.StatusUnauthorized, "Unauthorized", nil)
			}

			// Check if user role is in allowed roles
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return helpers.ResponseHttp(c, http.StatusForbidden, "Insufficient permissions", nil)
		}
	}
}
