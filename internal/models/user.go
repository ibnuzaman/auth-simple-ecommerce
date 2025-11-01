package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID                     int        `json:"-" gorm:"primaryKey"`
	Username               string     `json:"username" gorm:"column:username;type:varchar(20);not null;uniqueIndex:ux_users_username"`
	Email                  string     `json:"email" gorm:"column:email;type:varchar(100);not null;uniqueIndex:ux_users_email"`
	PhoneNumber            string     `json:"phone_number" gorm:"column:phone_number;type:varchar(15);not null;uniqueIndex:ux_users_phone"`
	FullName               string     `json:"full_name" gorm:"column:full_name;type:varchar(100);not null"`
	Address                string     `json:"address" gorm:"column:address;type:text"`
	Dob                    *time.Time `json:"dob,omitempty" gorm:"column:dob;type:date"`
	Password               string     `json:"-" gorm:"column:password;type:varchar(255);not null"`
	Role                   string     `json:"role,omitempty" gorm:"column:role;type:varchar(10);not null;default:'user'"`
	ResetPasswordToken     *string    `json:"-" gorm:"column:reset_password_token;type:varchar(255)"`
	ResetPasswordExpiry    *time.Time `json:"-" gorm:"column:reset_password_expiry;type:timestamp"`
	EmailVerificationToken *string    `json:"-" gorm:"column:email_verification_token;type:varchar(255)"`
	EmailVerified          bool       `json:"email_verified" gorm:"column:email_verified;default:false"`
	IsActive               bool       `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt              time.Time  `json:"-" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt              time.Time  `json:"-" gorm:"column:updated_at;autoUpdateTime"`
}

func (*User) TableName() string {
	return "users"
}

func (l User) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

type UserSession struct {
	ID                  int `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UserID              int       `json:"user_id" gorm:"type:int" validate:"required"`
	Token               string    `json:"token" gorm:"type:text" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"type:text" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
}

func (*UserSession) TableName() string {
	return "user_sessions"
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
