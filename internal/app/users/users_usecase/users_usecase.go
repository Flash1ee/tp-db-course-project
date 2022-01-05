package users_usecase

import (
	"github.com/jackc/pgx/v4"
	"tp-db-project/internal/app"
	"tp-db-project/internal/app/users"
	"tp-db-project/internal/app/users/models"
)

type UsersUsecase struct {
	repo users.Repository
}

func NewUsersUsecase(repo users.Repository) *UsersUsecase {
	return &UsersUsecase{
		repo: repo,
	}
}

func (u *UsersUsecase) CreateUser(user *models.User) (*models.User, error) {
	if err := u.repo.Create(user); err != nil {
		return nil, &app.GeneralError{
			Err:         ConstraintError,
			ExternalErr: err,
		}
	}
	return user, nil
}
func (u *UsersUsecase) GetUserByNickname(nickname string) (*models.User, error) {
	user, err := u.repo.GetByNickname(nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		}
		return nil, err
	}
	return user, nil
}
func (u *UsersUsecase) GetUser(nickname, email string) ([]*models.User, error) {
	someUsers, err := u.repo.GetByEmailOrNickname(nickname, email)
	if err != nil {
		return nil, &app.GeneralError{
			Err: err,
		}
	}
	return someUsers, nil
}
func (u *UsersUsecase) UpdateUser(user *models.User) (*models.User, error) {
	user, err := u.repo.Update(user)
	//email already exists
	if err != nil {
		return nil, &app.GeneralError{
			Err:         ConstraintError,
			ExternalErr: err,
		}
	}
	return user, nil
}
