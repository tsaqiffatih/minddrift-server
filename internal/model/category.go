package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"type:text"`
	Articles    []Article `gorm:"many2many:article_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
