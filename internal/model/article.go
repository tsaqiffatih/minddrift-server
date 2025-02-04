package model

import (
	"time"

	"github.com/google/uuid"
)

type ArticleStatus string

const (
	Draft     ArticleStatus = "draft"
	Review    ArticleStatus = "review"
	Published ArticleStatus = "published"
)

type Article struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title       string        `gorm:"not null"`
	Content     string        `gorm:"type:text;not null"`
	Slug        string        `gorm:"unique;not null"`
	Status      ArticleStatus `gorm:"type:varchar(10);default:'draft'"`
	AuthorID    uuid.UUID     `gorm:"type:uuid;not null"`
	PublishedAt *time.Time    `gorm:"default:null"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime"`

	Author   User       `gorm:"foreignKey:AuthorID"`
	Category []Category `gorm:"many2many:article_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []Tag      `gorm:"many2many:article_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
