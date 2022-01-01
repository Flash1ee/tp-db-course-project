package thread_repository

import (
	"database/sql"
	"tp-db-project/internal/app/thread/models"
)

const (
	queryGetThreadById = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE id = $1"
	queryGetThreadBySlug = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE slug = $1"
)

type ThreadRepository struct {
	conn *sql.DB
}

func NewThreadRepository(conn *sql.DB) *ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}

func (r *ThreadRepository) GetByID(id int64) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}

	if err := r.conn.QueryRow(queryGetThreadById, id).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *ThreadRepository) GetBySlug(slug string) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}

	if err := r.conn.QueryRow(queryGetThreadBySlug, slug).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		return nil, err
	}

	return res, nil
}
