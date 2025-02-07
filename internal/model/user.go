package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	Admin   UserRole = "admin"
	Editor  UserRole = "editor"
	Penulis UserRole = "penulis"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
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

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
