package thread_repository

import (
	"database/sql"
	"github.com/go-openapi/strfmt"
	"time"
	"tp-db-project/internal/app/thread/models"
	"tp-db-project/internal/app/thread/thread_repository"
)

const (
	queryGetThreadById = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE id = $1;"
	queryGetThreadBySlug = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE slug = $1;"

	queryUpdateThreadById = "UPDATE thread SET title = $2, message = $3 " +
		"WHERE id = $1 RETURNING id, title, author, forum, message, votes, slug, created;"
	queryUpdateThreadBySlug = "UPDATE thread SET title = $2, message = $3 " +
		"WHERE slug = $1 RETURNING id, title, author, forum, message, votes, slug, created;"
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
		if err == sql.ErrNoRows {
			return nil, thread_repository.NotFound
		}
		return nil, err
	}

	return res, nil
}

func (r *ThreadRepository) GetBySlug(slug string) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}

	if err := r.conn.QueryRow(queryGetThreadBySlug, slug).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		if err == sql.ErrNoRows {
			return nil, thread_repository.NotFound
		}

		return nil, err
	}

	return res, nil
}

func (r *ThreadRepository) UpdateByID(id int64, req *models.RequestUpdateThread) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	threadTime := &time.Time{}
	if err := r.conn.QueryRow(queryUpdateThreadById, id, req.Title, req.Message).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, threadTime); err != nil {
		return nil, err
	}

	res.Created = strfmt.DateTime(threadTime.UTC()).String()

	return res, nil
}

func (r *ThreadRepository) UpdateBySlug(slug string, req *models.RequestUpdateThread) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	threadTime := &time.Time{}
	if err := r.conn.QueryRow(queryUpdateThreadBySlug, slug, req.Title, req.Message).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, threadTime); err != nil {
		return nil, err
	}

	res.Created = strfmt.DateTime(threadTime.UTC()).String()

	return res, nil
}
