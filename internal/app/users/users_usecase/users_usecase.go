package users_usecase

import (
	"tp-db-project/internal/app/users"
)

type UsersUsecase struct {
	repo users.Repository
}

func NewUsersUsecase(repo users.Repository) *UsersUsecase {
	return &UsersUsecase{
		repo: repo,
	}
}
