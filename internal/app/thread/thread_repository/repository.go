package thread_repository

import "tp-db-project/internal/app/thread/models"

//go:generate mockgen -destination=mocks/repository.go -package=mock_thread -mock_names=Repository=ThreadRepository . Repository

type Repository interface {
	GetByID(id int64) (*models.ResponseThread, error)
	GetBySlug(slug string) (*models.ResponseThread, error)
	UpdateByID(id int64, req *models.RequestUpdateThread) (*models.ResponseThread, error)
	UpdateBySlug(slug string, req *models.RequestUpdateThread) (*models.ResponseThread, error)
}
