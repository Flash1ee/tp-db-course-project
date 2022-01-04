package forum

import (
	"tp-db-project/internal/app/forum/models"
	models_utilits "tp-db-project/internal/app/models"
	models_thread "tp-db-project/internal/app/thread/models"
	models_users "tp-db-project/internal/app/users/models"
)

//go:generate mockgen -destination=repository/mocks/repository.go -package=mock_forum -mock_names=Repository=ForumRepository . Repository

type Repository interface {
	Create(req *models.RequestCreateForum) error
	GetForumBySlag(slag string) (*models.Forum, error)
	GetForumUsers(slug string, since int, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error)
	GetForumThreads(slug string, since int, desc bool, pag *models_utilits.Pagination) ([]*models_thread.Thread, error)
}
