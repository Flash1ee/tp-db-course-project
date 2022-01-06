package thread

import (
	models2 "tp-db-project/internal/app/forum/models"
	pag_models "tp-db-project/internal/app/models"
	post_models "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/app/thread/models"
)

//go:generate mockgen -destination=mocks/repository.go -package=mock_thread -mock_names=Repository=ThreadRepository . Repository

type Repository interface {
	CreateThread(req *models2.RequestCreateThread) (*models.ResponseThread, error)
	GetByID(id int64) (*models.ResponseThread, error)
	GetBySlug(slug string) (*models.ResponseThread, error)
	UpdateByID(id int64, req *models.RequestUpdateThread) (*models.ResponseThread, error)
	UpdateBySlug(slug string, req *models.RequestUpdateThread) (*models.ResponseThread, error)
	CreatePosts(forum string, thread int64, posts []*models.RequestNewPost) ([]post_models.ResponsePost, error)
	GetPostsByFlats(id int, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error)
	GetPostsByTree(id int, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error)
	GetPostsByParentTree(id int, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error)
}
