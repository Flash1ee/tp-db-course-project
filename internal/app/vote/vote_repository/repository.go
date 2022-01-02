package vote_repository

import "tp-db-project/internal/app/vote/models"

//go:generate mockgen -destination=mocks/repository.go -package=mock_vote -mock_names=Repository=VoteRepository . Repository

type Repository interface {
	Create(vote *models.Vote) error
	Update(threadID int64, req *models.RequestVoteUpdate) (bool, error)
}
