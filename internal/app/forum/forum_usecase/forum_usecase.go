package forum_usecase

import (
	"tp-db-project/internal/app/forum"
)

type ForumUsecase struct {
	repo forum.Repository
}

func NewForumUsecase(repo forum.Repository) *ForumUsecase {
	return &ForumUsecase{
		repo: repo,
	}
}
