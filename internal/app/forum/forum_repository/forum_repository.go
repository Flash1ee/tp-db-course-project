package forum_repository

import (
	"database/sql"
	"fmt"
	"time"
	"tp-db-project/internal/app/forum/models"
	models_utilits "tp-db-project/internal/app/models"
	models_thread "tp-db-project/internal/app/thread/models"
	models_users "tp-db-project/internal/app/users/models"
)

const (
	queryCreate                    = "INSERT INTO forum(title, users_nickname, slug) VALUES($1, $2, $3);"
	queryGetForumBySlug            = "SELECT title, users_nickname, slug, posts, threads FROM forum WHERE slug = $1"
	queryCreateForumThreadWithTime = "INSERT INTO thread(title, author, forum, message, created) VALUES($1, $2, $3, $4, $5)"
	queryCreateForumThreadNoTime   = "INSERT INTO thread(title, author, forum, message) VALUES($1, $2, $3, $4)"
	queryGetForumUsers             = "SELECT nickname, fullname, about, email FROM forum " +
		"JOIN thread t ON forum.slug = t.forum JOIN post p ON t.id = p.thread " +
		"JOIN users u ON (p.author = u.nickname OR t.author = u.nickname) " +
		"WHERE forum.slug = $1 "
	queryGetForumThreads = "SELECT t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM threads as t " +
		"LEFT JOIN forum f on t.forum = f.slug " +
		"WHERE f.slug = $1 "
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
func (r *ForumRepository) CreateForumThread(forumName string, req *models.RequestCreateThread) error {
	var err error
	if req.Created == "" {
		_, err = r.conn.Exec(queryCreateForumThreadNoTime, req.Title, req.Author, forumName, req.Message)
	} else {
		_, err = r.conn.Exec(queryCreateForumThreadWithTime, req.Title, req.Author, forumName, req.Message, req.Created)
	}

	return err
}
func (r *ForumRepository) GetForumUsers(slug string, since string, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error) {
	orderBy := "ORDER BY u.nickname "
	querySince := "AND u.nickname > $2"
	query := queryGetForumUsers
	limit := pag.Limit
	var rows *sql.Rows
	var err error

	if desc {
		orderBy += "DESC"
	}
	if limit > 0 {
		orderBy += fmt.Sprintf("LIMIT %d", pag.Limit)
	}
	if since != "" {
		query = query + querySince + orderBy
		rows, err = r.conn.Query(query, slug, since)
	} else {
		query = query + orderBy
		rows, err = r.conn.Query(query, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models_users.User
	for rows.Next() {
		user := &models_users.User{}
		if err := rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
			return nil, err
		}
		res = append(res, user)
	}

	return res, nil
}
func (r *ForumRepository) GetForumThreads(slug string, since string, desc bool, pag *models_utilits.Pagination) ([]*models_thread.Thread, error) {
	orderBy := "ORDER BY t.created "
	querySince := "AND t.created > $2"
	query := queryGetForumThreads
	limit := pag.Limit

	var rows *sql.Rows
	var err error

	if desc {
		orderBy += "DESC"
	}
	if limit > 0 {
		orderBy += fmt.Sprintf("LIMIT %d", pag.Limit)
	}
	if since != "" {
		query = query + querySince + orderBy
		rows, err = r.conn.Query(query, slug, since)
	} else {
		query = query + orderBy
		rows, err = r.conn.Query(query, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models_thread.Thread
	for rows.Next() {
		thread := &models_thread.Thread{}
		createdTime := time.Time{}
		if err := rows.Scan(&thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &createdTime); err != nil {
			return nil, err
		}
		res = append(res, thread)
	}

	return res, nil
}