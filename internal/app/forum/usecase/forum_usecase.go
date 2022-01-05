package forum_usecase

import (
	"github.com/jackc/pgx/v4"
	"tp-db-project/internal/app"
	"tp-db-project/internal/app/forum"
	"tp-db-project/internal/app/forum/models"
	models_utilits "tp-db-project/internal/app/models"
	"tp-db-project/internal/app/thread"
	models_thread "tp-db-project/internal/app/thread/models"
	"tp-db-project/internal/app/users"
	models_users "tp-db-project/internal/app/users/models"
	"tp-db-project/internal/app/users/users_usecase"
)

type ForumUsecase struct {
	repo       forum.Repository
	usersRepo  users.Repository
	threadRepo thread.Repository
}

func NewForumUsecase(repo forum.Repository, usersRepo users.Repository, thRepo thread.Repository) *ForumUsecase {
	return &ForumUsecase{
		repo:       repo,
		usersRepo:  usersRepo,
		threadRepo: thRepo,
	}
}

func (u *ForumUsecase) Create(req *models.RequestCreateForum) (*models.Forum, error) {
	if f, err := u.repo.GetForumBySlag(req.Slug); err == nil {
		return f, AlreadyExists
	}
	if _, err := u.usersRepo.Get(req.User); err != nil {
		if err == pgx.ErrNoRows {
			return nil, users_usecase.NotFound
		}
		return nil, err
	}

	if err := u.repo.Create(req); err != nil {
		return nil, app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}
	return nil, nil
}
func (u *ForumUsecase) GetForum(slug string) (*models.Forum, error) {
	res, err := u.repo.GetForumBySlag(slug)
	if err != nil {
		return nil, app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	}
	return res, err
}
func (u *ForumUsecase) GetForumUsers(slug string, since int, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error) {
	if _, err := u.repo.GetForumBySlag(slug); err != nil {
		return nil, app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	}
	res, err := u.repo.GetForumUsers(slug, since, desc, pag)
	if err != nil {
		return nil, app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}
	return res, nil
}
func (u *ForumUsecase) GetForumThreads(slug string, since int, desc bool, pag *models_utilits.Pagination) ([]*models_thread.Thread, error) {
	if _, err := u.repo.GetForumBySlag(slug); err != nil {
		return nil, app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	}
	res, err := u.repo.GetForumThreads(slug, since, desc, pag)
	if err != nil {
		return nil, app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}
	return res, nil
}
func (u *ForumUsecase) CreateThread(forumName string, req *models.RequestCreateThread) (*models_thread.ResponseThread, error) {
	if _, err := u.usersRepo.Get(req.Author); err != nil {
		if err == pgx.ErrNoRows {
			return nil, users_usecase.NotFound
		}
		return nil, err
	}
	if _, err := u.repo.GetForumBySlag(forumName); err == nil {
		return nil, AlreadyExists
	}

	th, err := u.threadRepo.GetBySlug(forumName)
	if err == nil {
		return th, AlreadyExists
	}
	th, err = u.threadRepo.CreateThread(forumName, req)
	if err != nil {
		return nil, app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	} else {
		return th, nil
	}
}
