package repository

import "github.com/tsaqiffatih/minddrift-server/internal/model"

type ArticleVersionRepository interface {
	CreateArticleVersion(articleVersion model.ArticleVersion) (model.ArticleVersion, error)
}
