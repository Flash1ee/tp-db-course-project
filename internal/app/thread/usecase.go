package thread

import (
	"tp-db-project/internal/app/thread/models"
	models2 "tp-db-project/internal/app/vote/models"
)

//go:generate mockgen -destination=usecase/mocks/usecase.go -package=mock_thread -mock_names=Usecase=ThreadUsecase . Usecase

type Usecase interface {
	GetThreadInfo(slugOrID string) (*models.ResponseThread, error)
	UpdateThread(slugOrID string, req *models.RequestUpdateThread) (*models.ResponseThread, error)
	UpdateVoice(slugOrID string, req *models2.RequestVoteUpdate) (bool, error)
}
