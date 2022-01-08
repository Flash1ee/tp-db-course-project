package vote_postgresql

import (
	"github.com/jackc/pgx"
	"tp-db-project/internal/app/vote/models"
	"tp-db-project/internal/app/vote/repository"
)

const (
	queryCreateVote  = "INSERT INTO vote(nickname, thread_id, voice) VALUES($1, $2, $3);"
	queryUpdateVote  = "UPDATE vote SET voice = $3 WHERE thread_id = $1 and nickname = $2 and voice != $3;"
	queryCheckExists = "SELECT voice from vote where nickname = $1 and thread_id = $2;"
)

type VoteRepository struct {
	conn *pgx.ConnPool
}

func NewVoteRepository(conn *pgx.ConnPool) *VoteRepository {
	return &VoteRepository{
		conn: conn,
	}
}
func (r *VoteRepository) Exists(nickname string, threadID int64) (bool, error) {
	voice := 0
	err := r.conn.QueryRow(queryCheckExists, nickname, threadID).
		Scan(&voice)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (r *VoteRepository) Create(vote *models.Vote) error {
	if vote == nil {
		return repository.InvalidArgument
	}
	if _, err := r.conn.Exec(queryCreateVote, vote.Nickname, vote.ThreadID, vote.Voice); err != nil {
		return err
	}
	return nil
}
func (r *VoteRepository) Update(threadID int64, req *models.RequestVoteUpdate) (bool, error) {
	if req == nil {
		return false, repository.InvalidArgument
	}
	res, err := r.conn.Exec(queryUpdateVote, threadID, req.Nickname, req.Voice)
	if err != nil {
		return false, err
	}
	UpdatedRows := res.RowsAffected()
	return UpdatedRows == 1, nil
}
