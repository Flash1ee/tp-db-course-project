package forum_postgresql

import (
	"fmt"
	"github.com/jackc/pgx"
	"tp-db-project/internal/app/forum/models"
	models_utilits "tp-db-project/internal/app/models"
	models_thread "tp-db-project/internal/app/thread/models"
	models_users "tp-db-project/internal/app/users/models"
)

const (
	queryGetForumUsers = "SELECT DISTINCT u.nickname, u.fullname, u.about, u.email from user_forum " +
		"LEFT JOIN users u on user_forum.nickname = u.nickname " +
		"where user_forum.forum = $1 "
	queryCreate         = "INSERT INTO forum(title, users_nickname, slug) VALUES($1, $2, $3);"
	queryGetForumBySlug = "SELECT title, users_nickname, slug, posts, threads FROM forum WHERE slug = $1"

	//queryGetForumUsers = "SELECT DISTINCT u.nickname, u.fullname, u.about, u.email FROM forum " +
	//	"JOIN thread t ON forum.slug = t.forum JOIN post p ON t.id = p.thread " +
	//	"JOIN users u ON (p.author = u.nickname OR t.author = u.nickname) " +
	//	"WHERE forum.slug = $1 "
	queryGetForumThreads = "SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM thread as t " +
		"LEFT JOIN forum f on t.forum = f.slug " +
		"WHERE f.slug = $1 "
)

type ForumRepository struct {
	conn *pgx.ConnPool
}

func NewForumRepository(conn *pgx.ConnPool) *ForumRepository {
	return &ForumRepository{
		conn: conn,
	}
}

func (r *ForumRepository) Create(req *models.RequestCreateForum) error {
	_, err := r.conn.Exec(queryCreate, &req.Title, &req.User, &req.Slug)
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

func (r *ForumRepository) GetForumUsers(slug string, since string, desc bool, pag *models_utilits.Pagination) ([]*models_users.User, error) {
	//orderBy := "ORDER BY u.nickname "
	//querySince := " AND u.nickname > $2 "
	query := queryGetForumUsers
	limit := pag.Limit
	var rows *pgx.Rows
	var err error

	if desc && since != "" {
		query += fmt.Sprintf(" and u.nickname < '%s'", since)
	} else if since != "" {
		query += fmt.Sprintf(" and u.nickname > '%s'", since)
	}
	query += " ORDER BY u.nickname "
	if desc {
		query += "desc"
	}
	query += fmt.Sprintf(" LIMIT %d", limit)

	//if desc {
	//	orderBy += "DESC"
	//}
	//if limit > 0 {
	//	orderBy += fmt.Sprintf(" LIMIT %d", pag.Limit)
	//}
	//if since != -1 {
	//	//query = query + querySince + orderBy
	//	rows, err = r.conn.Query(context.Background(), query, slug, since)
	//} else {
	//query = query + orderBy
	rows, err = r.conn.Query(query, slug)
	//}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models_users.User, 0, 0)
	for rows.Next() {
		user := &models_users.User{}
		if err := rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
			return nil, err
		}
		res = append(res, user)
	}

	return res, nil
}
func (r *ForumRepository) GetForumThreads(forumSlug string, sinceDate string, desc bool, pag *models_utilits.Pagination) ([]*models_thread.ResponseThread, error) {
	orderBy := "ORDER BY t.created "
	querySince := " AND t.created >= $2 "
	query := queryGetForumThreads
	limit := pag.Limit

	var rows *pgx.Rows
	var err error

	if desc {
		orderBy += "DESC"
	}
	if limit > 0 {
		orderBy += fmt.Sprintf(" LIMIT %d", pag.Limit)
	}

	if sinceDate != "" && desc {
		querySince = " and t.created <= $2 "
	} else if sinceDate != "" && !desc {
		querySince = " and t.created >= $2 "
	}

	if sinceDate != "" {
		query = query + querySince + orderBy
		rows, err = r.conn.Query(query, forumSlug, sinceDate)
	} else {
		query = query + orderBy
		rows, err = r.conn.Query(query, forumSlug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*models_thread.ResponseThread, 0, 0)
	for rows.Next() {
		thread := &models_thread.ResponseThread{}
		//createdTime := time.Time{}
		//createdTime := &strfmt.DateTime{}
		if err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
			return nil, err
		}
		//thread.Created = strfmt.DateTime(createdTime.UTC()).String()
		//thread.Created = time.Time(*createdTime).UTC()
		res = append(res, thread)
	}

	return res, nil
}
