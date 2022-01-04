package vote

import "tp-db-project/internal/app/vote/models"

//go:generate mockgen -destination=mocks/repository.go -package=mock_vote -mock_names=Repository=VoteRepository . Repository

type Repository interface {
	Exists(nickname string, threadID int64) (bool, error)
	Create(vote *models.Vote) error
	Update(threadID int64, req *models.RequestVoteUpdate) (bool, error)
}
