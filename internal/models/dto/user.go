package dto

import "time"

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=20" example:"admin"`
	Email       string `json:"email" validate:"required,email,max=100" example:"admin@example.com"`
	PhoneNumber string `json:"phone_number" validate:"required,min=8,max=15" example:"62877618152"`
	FullName    string `json:"full_name" validate:"required,min=3,max=100" example:"super admin"`
	Address     string `json:"address" validate:"omitempty,max=500" example:"Jl Patiunus 1"`
	Dob         string `json:"dob" validate:"omitempty,datetime=2006-01-02" example:"1999-01-01"` // "YYYY-MM-DD"
	Password    string `json:"password" validate:"required,min=8,max=72" example:"password"`
}

type RegisterResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	FullName    string    `json:"full_name"`
	Address     string    `json:"address,omitempty"`
	Dob         string    `json:"dob,omitempty"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}
