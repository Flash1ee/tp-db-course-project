package users_usecase

import (
	"database/sql"
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
	if _, err := u.repo.Get(user.Nickname); err != sql.ErrNoRows {
		return nil, app.GeneralError{
			Err:         AlreadyExistsErr,
			ExternalErr: err,
		}
	}
	if err := u.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UsersUsecase) GetUser(nickname string) (*models.User, error) {
	user, err := u.repo.Get(nickname)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UsersUsecase) UpdateUser(user *models.User) (*models.User, error) {
	_, err := u.repo.Get(user.Nickname)
	if err != nil {
		return nil, app.GeneralError{
			Err:         NotFound,
			ExternalErr: err,
		}
	}
	user, err = u.repo.Update(user)
	//email already exists
	if err != nil {
		return nil, app.GeneralError{
			Err:         ConstraintError,
			ExternalErr: err,
		}
	}
	return user, nil
}
