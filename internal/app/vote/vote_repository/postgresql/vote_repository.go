package vote_postgresql

import (
	"database/sql"
	"tp-db-project/internal/app/vote/models"
	"tp-db-project/internal/app/vote/vote_repository"
)

const (
	queryCreateVote = "INSERT INTO vote(nickname, thread_id, voice) VALUES($1, $2, $3);"
	queryUpdateVote = "UPDATE vote SET voice = $3 WHERE thread_id = $1 and nickname = $2 and voice != $3;"
)

type VoteRepository struct {
	conn *sql.DB
}

func NewVoteRepository(conn *sql.DB) *VoteRepository {
	return &VoteRepository{
		conn: conn,
	}
}

func (r *VoteRepository) Create(vote *models.Vote) error {
	if vote == nil {
		return vote_repository.InvalidArgument
	}
	if _, err := r.conn.Exec(queryCreateVote, vote.Nickname, vote.ThreadID, vote.Voice); err != nil {
		return err
	}
	return nil
}
func (r *VoteRepository) Update(threadID int64, req *models.RequestVoteUpdate) (bool, error) {
	if req == nil {
		return false, vote_repository.InvalidArgument
	}
	res, err := r.conn.Exec(queryUpdateVote, threadID, req.Nickname, req.Voice)
	if err != nil {
		return false, err
	}
	UpdatedRows, _ := res.RowsAffected()
	return UpdatedRows == 1, nil
}
