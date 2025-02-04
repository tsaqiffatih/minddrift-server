package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Analytic struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ArticleID       uuid.UUID `gorm:"type:uuid;not null"`
	Views           int       `gorm:"default:0"`
	UniqueVisitors  int       `gorm:"default:0"`
	AverageTimeSpent int       `gorm:"default:0"`
	Article         Article   `gorm:"foreignKey:ArticleID"`
}
