package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ibnuzaman/auth-simple-ecommerce.git/internal/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser creates a new user
func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail finds user by email
func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone
func (r *AuthRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("phone_number = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

// FindByUsername finds user by username
func (r *AuthRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailOrUsername finds user by email or username
func (r *AuthRepository) FindByEmailOrUsername(ctx context.Context, emailOrUsername string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds user by ID
func (r *AuthRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdatePassword updates user password
func (r *AuthRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error
}

// SaveResetToken saves password reset token
func (r *AuthRepository) SaveResetToken(ctx context.Context, userID int, token string, expiry string) error {
	expiryTime, err := time.Parse(time.RFC3339, expiry)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"reset_password_token":  token,
			"reset_password_expiry": expiryTime,
		}).Error
}

// FindByResetToken finds user by reset token
func (r *AuthRepository) FindByResetToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("reset_password_token = ? AND reset_password_expiry > ?", token, time.Now()).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// ClearResetToken clears password reset token
func (r *AuthRepository) ClearResetToken(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"reset_password_token":  nil,
			"reset_password_expiry": nil,
		}).Error
}

// CreateSession creates a new user session
func (r *AuthRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// FindSessionByToken finds session by token
func (r *AuthRepository) FindSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.WithContext(ctx).
		Where("token = ? AND token_expired > ?", token, time.Now()).
		First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// FindSessionByUserID finds session by user ID
func (r *AuthRepository) FindSessionByUserID(ctx context.Context, userID int) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// UpdateSession updates user session
func (r *AuthRepository) UpdateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// DeleteSession deletes session by token
func (r *AuthRepository) DeleteSession(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.UserSession{}).Error
}

// DeleteSessionsByUserID deletes all sessions for a user
func (r *AuthRepository) DeleteSessionsByUserID(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}
