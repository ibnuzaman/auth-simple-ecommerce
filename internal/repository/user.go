package repository

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/constants"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

var reConstraint = regexp.MustCompile(`constraint\s+"([^"]+)"`)

func mapConstraintToField(constraint string) string {
	c := strings.ToLower(constraint)
	switch {
	case strings.Contains(c, "ux_users_email"):
		return "email"
	case strings.Contains(c, "ux_users_username"):
		return "username"
	case strings.Contains(c, "ux_users_phone"):
		return "phone_number"
	default:
		return "unknown"
	}
}

func (r *UserRepository) InsertNewUser(ctx context.Context, user *models.User) error {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// Coba tebak field dari string error (karena gorm.ErrDuplicatedKey tidak bawa detail)
			field := "unknown"
			if m := reConstraint.FindStringSubmatch(err.Error()); len(m) == 2 {
				field = mapConstraintToField(m[1])
			}
			return &constants.ConflictError{Field: field}
		}

		return err
	}

	return nil
}
