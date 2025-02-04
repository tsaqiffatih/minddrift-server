package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SEOMetadata struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ArticleID       uuid.UUID `gorm:"type:uuid;not null"`
	MetaTitle       string
	MetaDescription string  `gorm:"type:text"`
	Keywords        string  `gorm:"type:text"`
	Article         Article `gorm:"foreignKey:ArticleID"`
}
