package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ArticleID uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Content   string    `gorm:"type:text;not null"`
	Article   Article   `gorm:"foreignKey:ArticleID"`
	User      User      `gorm:"foreignKey:UserID"`
}
