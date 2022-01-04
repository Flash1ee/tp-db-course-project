package thread

import (
	pag_models "tp-db-project/internal/app/models"
	post_models "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/app/thread/models"
	models2 "tp-db-project/internal/app/vote/models"
)

//go:generate mockgen -destination=usecase/mocks/usecase.go -package=mock_thread -mock_names=Usecase=ThreadUsecase . Usecase

type Usecase interface {
	GetThreadInfo(slugOrID string) (*models.ResponseThread, error)
	UpdateThread(slugOrID string, req *models.RequestUpdateThread) (*models.ResponseThread, error)
	UpdateVoice(slugOrID string, req *models2.RequestVoteUpdate) (bool, error)
	GetPostsBySort(slugOrId string, sort string, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error)
}
