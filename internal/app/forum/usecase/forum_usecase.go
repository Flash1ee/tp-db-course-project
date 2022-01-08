package forum_usecase

import (
	"github.com/jackc/pgx"
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
	f, err := u.repo.GetForumBySlag(req.Slug)
	if err == nil {
		return f, AlreadyExists
	}
	author, err := u.usersRepo.GetByNickname(req.User)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &app.GeneralError{
				Err:         users_usecase.NotFound,
				ExternalErr: err,
			}
		}
		return nil, err
	}
	req.User = author.Nickname

	if err := u.repo.Create(req); err != nil {
		return nil, &app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}
	return &models.Forum{
		Title:         req.Title,
		Slug:          req.Slug,
		UsersNickname: req.User,
	}, nil
}
func (u *ForumUsecase) GetForum(slug string) (*models.Forum, error) {
	res, err := u.repo.GetForumBySlag(slug)
	if err != nil {
		return nil, &app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	}
	return res, err
}
func (u *ForumUsecase) GetForumUsers(slug string, since string, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error) {
	if _, err := u.repo.GetForumBySlag(slug); err != nil {
		return nil, &app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	}
	res, err := u.repo.GetForumUsers(slug, since, desc, pag)
	if err != nil {
		return nil, &app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}
	return res, nil
}
func (u *ForumUsecase) GetForumThreads(forumSlug string, sinceDate string, desc bool, pag *models_utilits.Pagination) ([]*models_thread.ResponseThread, error) {
	if _, err := u.repo.GetForumBySlag(forumSlug); err != nil {
		return nil, &app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	}
	res, err := u.repo.GetForumThreads(forumSlug, sinceDate, desc, pag)
	if err != nil {
		return nil, &app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}
	return res, nil
}
func (u *ForumUsecase) CreateThread(req *models.RequestCreateThread) (*models_thread.ResponseThread, error) {
	if _, err := u.usersRepo.GetByNickname(req.Author); err != nil {
		if err == pgx.ErrNoRows {
			return nil, users_usecase.NotFound
		}
		return nil, err
	}
	if f, err := u.repo.GetForumBySlag(req.Forum); err != nil {
		return nil, &app.GeneralError{
			Err:         ForumNotFound,
			ExternalErr: err,
		}
	} else {
		req.Forum = f.Slug
	}
	if req.Slug != "" {
		if th, err := u.threadRepo.GetBySlug(req.Slug); err == nil {
			//if th.Forum == req.Forum && th.Slug == req.Slug {
			return th, AlreadyExists
			//}
		}
	}

	th, err := u.threadRepo.CreateThread(req)
	if err != nil {
		return nil, &app.GeneralError{
			Err:         ConstraintError,
			ExternalErr: err,
		}
	} else {
		return th, nil
	}
}
