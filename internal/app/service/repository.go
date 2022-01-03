package service

import "tp-db-project/internal/app/service/models"

//go:generate mockgen -destination=mocks/repository.go -package=mock_service -mock_names=Repository=ServiceRepository . Repository

type Repository interface {
	Delete() error
	Status() (*models.ResponseService, error)
}
