package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Backup struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FilePath  string    `gorm:"not null"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:CreatedBy"`
}
