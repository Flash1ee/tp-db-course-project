package post_repository

import "database/sql"

type PostRepository struct {
	conn *sql.DB
}

func NewPostRepository(conn *sql.DB) *PostRepository {
	return &PostRepository{
		conn: conn,
	}
}
