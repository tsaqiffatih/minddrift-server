package repository

import (
	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
)

type CommentRepository interface {
	CreateComment(comment *model.Comment) error
	GetCommentByID(id uuid.UUID) (*model.Comment, error)
	GetCommentByArticle(articleID uuid.UUID) ([]model.Comment, error)
	DeleteComent(id uuid.UUID) error
}
