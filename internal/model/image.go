package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	URL        string    `gorm:"not null"`
	AltText    string
	Caption    string
	ArticleID  uuid.UUID `gorm:"type:uuid;not null"`
	Article    Article   `gorm:"foreignKey:ArticleID"`
	UploadedBy uuid.UUID `gorm:"type:uuid;not null"`
	User       User      `gorm:"foreignKey:UploadedBy"`
}
