package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRole string

const (
	Admin   UserRole = "admin"
	Editor  UserRole = "editor"
	Penulis UserRole = "penulis"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username      string    `gorm:"unique;not null" json:"username"`
	Email         string    `gorm:"unique;not null" json:"email"`
	Password      string    `gorm:"not null" json:"-"`
	Role          UserRole  `gorm:"type:varchar(10);default:'penulis'" json:"role"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	TwoFAEnabled  bool      `gorm:"default:false" json:"two_fa_enabled"`
	TwoFASecret   string    `gorm:"type:text" json:"-"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
