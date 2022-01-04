package post

import "tp-db-project/internal/app/post/models"

//go:generate mockgen -destination=repository/mocks/repository.go -package=mock_post -mock_names=Repository=PostRepository . Repository

type Repository interface {
	Get(id int64, related string) (*models.Post, error)
	Update(id int64, req *models.RequestUpdateMessage) (*models.ResponsePost, error)
}
