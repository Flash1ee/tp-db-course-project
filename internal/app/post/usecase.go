package post

import "tp-db-project/internal/app/post/models"

//go:generate mockgen -destination=usecase/mocks/usecase.go -package=mock_post -mock_names=Usecase=PostUsecase . Usecase

type Usecase interface {
	GetPost(id int64, related string) (*models.ResponsePostDetail, error)
	UpdatePost(id int64, req *models.RequestUpdateMessage) (*models.ResponsePost, error)
}
