package users

import "tp-db-project/internal/app/users/models"

//go:generate mockgen -destination=usecase/mocks/usecase.go -package=mock_users -mock_names=Usecase=UsersUsecase . Usecase

type Usecase interface {
	GetUserByNickname(nickname string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	GetUser(nickname string, email string) ([]*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
}
