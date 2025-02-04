package repository

import (
	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
)

type ArticleRepository interface {
	CreateArticle(article model.Article) error
	GetArticleByID(id string) (model.Article, error)
	GetArticleByAuthor(authorID uuid.UUID) ([]model.Article, error)
	GetArticleBySlug(slug string) ([]model.Article, error)
	UpdateArticle(article model.Article) error
	ListAllArticle(status string, limit, offset int) ([]model.Article, error)
	DeleteArticle(id string) error
}
