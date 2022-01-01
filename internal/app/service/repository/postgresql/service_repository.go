package postgresql

import (
	"database/sql"
	"tp-db-project/internal/app/service/models"
)

const (
	queryDeleteAllTables           = "TRUNCATE TABLE forum, post, thread, user_forum, users, vote CASCADE;"
	queryCountForumPostThreadUsers = "SELECT (SELECT count(*) FROM forum) AS forum," +
		"(SELECT count(*) FROM post) AS post, " +
		"(SELECT count(*) FROM thread) AS thread, " +
		"(SELECT count(*) FROM users) AS users;"
)

type ServiceRepository struct {
	conn *sql.DB
}

func NewServiceRepository(conn *sql.DB) *ServiceRepository {
	return &ServiceRepository{
		conn: conn,
	}
}

func (r *ServiceRepository) Delete() error {
	_, err := r.conn.Exec(queryDeleteAllTables)
	return err
}

func (r *ServiceRepository) Status() (*models.ResponseService, error) {
	var res *models.ResponseService

	if err := r.conn.QueryRow(queryCountForumPostThreadUsers).
		Scan(&res.Forum, &res.Post, &res.Thread, &res.User); err != nil {
		return nil, err
	}

	return res, nil
}