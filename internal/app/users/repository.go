package users

import "tp-db-project/internal/app/users/models"

//go:generate mockgen -destination=repository/mocks/repository.go -package=mock_users -mock_names=Repository=UsersRepository . Repository

type Repository interface {
	GetByNickname(nickname string) (*models.User, error)
	Create(user *models.User) error
	GetByEmailOrNickname(nickname, email string) ([]*models.User, error)
	Update(user *models.User) (*models.User, error)
}
