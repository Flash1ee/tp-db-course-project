package forum_repository

import (
	"database/sql"
	"tp-db-project/internal/app/forum/models"
	models_utilits "tp-db-project/internal/app/models"
	models_users "tp-db-project/internal/app/users/models"
)

const (
	queryCreate         = "INSERT INTO forum(title, users_nickname, slug) VALUES($1, $2, $3);"
	queryGetForumBySlug = "SELECT title, users_nickname, slug, posts, threads FROM forum WHERE slug = $1"
)

type ForumRepository struct {
	conn *sql.DB
}

func NewForumRepository(conn *sql.DB) *ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}

func (r *ForumRepository) Create(req *models.RequestCreateForum) error {
	_, err := r.conn.Exec(queryCreate, req.Title, req.User, req.Slug)
	return err
}
func (r *ForumRepository) GetForumBySlag(slag string) (*models.Forum, error) {
	forum := &models.Forum{}

	res := r.conn.QueryRow(queryGetForumBySlug, slag).
		Scan(&forum.Title, &forum.UsersNickname, &forum.Slug, &forum.Posts, &forum.Threads)
	if res != nil {
		return nil, res
	}
	return forum, nil
}
func (r *ForumRepository) GetForumUsers(slug string, since int, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error) {
	//limit, _, err := utilits.AddPagination("forum", pag, r.conn)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}
