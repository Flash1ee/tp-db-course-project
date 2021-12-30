package thread_repository

import "database/sql"

type ThreadRepository struct {
	conn *sql.DB
}

func NewThreadRepository(conn *sql.DB) *ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}
