package service_usecase

import (
	"tp-db-project/internal/app/service"
	"tp-db-project/internal/app/service/models"
)

type ServiceUsecase struct {
	repo service.Repository
}

func NewServiceUsecase(repo service.Repository) *ServiceUsecase {
	return &ServiceUsecase{
		repo: repo,
	}
}

func (u *ServiceUsecase) Clear() error {
	return u.repo.Delete()
}

func (u *ServiceUsecase) Status() (*models.ResponseService, error) {
	return u.repo.Status()
}
