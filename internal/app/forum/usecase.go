package forum

import (
	"tp-db-project/internal/app/forum/models"
	models_utilits "tp-db-project/internal/app/models"
	models_thread "tp-db-project/internal/app/thread/models"
	models_users "tp-db-project/internal/app/users/models"
)

//go:generate mockgen -destination=usecase/mocks/usecase.go -package=mock_forum -mock_names=Usecase=ForumUsecase . Usecase

type Usecase interface {
	GetForum(slug string) (*models.Forum, error)
	Create(req *models.RequestCreateForum) (*models.Forum, error)
	CreateThread(req *models.RequestCreateThread) (*models_thread.ResponseThread, error)
	GetForumUsers(slug string, since string, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error)
	GetForumThreads(forumSlug string, sinceDate string, desc bool, pag *models_utilits.Pagination) ([]*models_thread.ResponseThread, error)
}
