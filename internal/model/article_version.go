package model

import (
	"time"

	"github.com/google/uuid"
)

type ArticleVersion struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ArticleID uuid.UUID `gorm:"type:uuid;not null"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Article Article `gorm:"foreignKey:ArticleID"`
}
