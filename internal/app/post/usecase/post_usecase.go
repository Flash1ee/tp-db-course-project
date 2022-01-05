package post_usecase

import (
	"tp-db-project/internal/app/post"
	"tp-db-project/internal/app/post/models"
	post_repository "tp-db-project/internal/app/post/repository"
)

type PostUsecase struct {
	repo post.Repository
}

func NewPostUsecase(repo post.Repository) *PostUsecase {
	return &PostUsecase{
		repo: repo,
	}
}

func (u *PostUsecase) GetPost(id int64, related string) (*models.ResponsePostDetail, error) {
	return u.repo.Get(id, related)
}
func (u *PostUsecase) UpdatePost(id int64, req *models.RequestUpdateMessage) (*models.ResponsePost, error) {
	if _, err := u.repo.Get(id, ""); err != nil {
		return nil, post_repository.NotFound
	}
	return u.repo.Update(id, req)
}
