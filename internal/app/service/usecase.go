package service

import "tp-db-project/internal/app/service/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mock_users -mock_names=Usecase=UsersUsecase . Usecase

type Usecase interface {
	Clear() error
	Status() (*models.ResponseService, error)
}
