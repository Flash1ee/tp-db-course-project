package forum_repository

import (
	"database/sql"
)

type ForumRepository struct {
	conn *sql.DB
}

func NewForumRepository(conn *sql.DB) *ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}
