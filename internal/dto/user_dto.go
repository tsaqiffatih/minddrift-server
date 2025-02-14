package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email,omitempty"`
	Password string `json:"password" validate:"required,min=8,max=32,omitempty"`
	Role     string `json:"role" validate:"oneof=admin editor penulis"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	NewPassword string `json:"password" validate:"required"`
	Token       string `json:"token" validate:"required"`
}

type ValidateResetToken struct {
	Token string `json:"token" validate:"required"`
}

type UpdateUserRequest struct {
	Username     *string `json:"username,omitempty" validate:"omitempty,min=3,max=20"`
	Email        *string `json:"email,omitempty" validate:"omitempty,email"`
	Password     *string `json:"password,omitempty" validate:"omitempty,min=8,max=32"`
	Role         *string `json:"role,omitempty" validate:"omitempty,oneof=admin editor penulis"`
	TwoFAEnabled *bool   `json:"two_fa_enabled,omitempty"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	TwoFAEnabled  bool      `json:"two_fa_enabled"`
}

type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         UserResponse `json:"user"`
}
