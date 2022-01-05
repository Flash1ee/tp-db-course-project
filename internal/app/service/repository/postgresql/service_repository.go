package service_postgresql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
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
	conn *pgxpool.Pool
}

func NewServiceRepository(conn *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		conn: conn,
	}
}

func (r *ServiceRepository) Delete() error {
	_, err := r.conn.Exec(context.Background(), queryDeleteAllTables)
	return err
}

func (r *ServiceRepository) Status() (*models.ResponseService, error) {
	res := &models.ResponseService{}

	if err := r.conn.QueryRow(context.Background(), queryCountForumPostThreadUsers).
		Scan(&res.Forum, &res.Post, &res.Thread, &res.User); err != nil {
		return nil, err
	}

	return res, nil
}
