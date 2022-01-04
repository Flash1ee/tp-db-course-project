package thread_repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"time"
	models2 "tp-db-project/internal/app/forum/models"
	pag_models "tp-db-project/internal/app/models"
	post_models "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/app/thread/models"
	"tp-db-project/internal/app/thread/thread_repository"
)

const (
	queryGetThreadById = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE id = $1;"
	queryGetThreadBySlug = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE slug = $1;"
	queryCreateForumThreadWithTime = "INSERT INTO thread(title, author, forum, message, created) VALUES($1, $2, $3, $4, $5) " +
		"RETURNING id, title, author, forum, message, votes, slug, created;"
	queryCreateForumThreadNoTime = "INSERT INTO thread(title, author, forum, message) VALUES($1, $2, $3, $4) " +
		"RETURNING id, title, author, forum, message, votes, slug, created; "
	queryUpdateThreadById = "UPDATE thread SET title = $2, message = $3 " +
		"WHERE id = $1 RETURNING id, title, author, forum, message, votes, slug, created;"
	queryUpdateThreadBySlug = "UPDATE thread SET title = $2, message = $3 " +
		"WHERE slug = $1 RETURNING id, title, author, forum, message, votes, slug, created;"

	queryInsertPost = "INSERT INTO post(parent, author, message, forum, thread, created) VALUES "
	queryGetPosts   = "SELECT id, parent, author, message, is_edited, forum, thread, created " +
		"FROM post WHERE thread = $1 "
)

type ThreadRepository struct {
	conn *pgx.Conn
}

func NewThreadRepository(conn *pgx.Conn) *ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}

func (r *ThreadRepository) CreateThread(forumName string, req *models2.RequestCreateThread) (*models.ResponseThread, error) {
	var err error
	res := &models.ResponseThread{}
	if req == nil {
		return nil, thread_repository.ArgError
	}
	threadTime := &time.Time{}

	if req.Created == "" {
		err = r.conn.QueryRow(context.Background(), queryCreateForumThreadNoTime, req.Title, req.Author, forumName, req.Message).
			Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
				&res.Message, &res.Votes, &res.Slug, threadTime)
	} else {
		err = r.conn.QueryRow(context.Background(), queryCreateForumThreadWithTime, req.Title, req.Author, forumName, req.Message, req.Created).
			Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
				&res.Message, &res.Votes, &res.Slug, threadTime)
	}

	res.Created = strfmt.DateTime(threadTime.UTC()).String()

	return res, err
}

func (r *ThreadRepository) GetByID(id int64) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	if err := r.conn.QueryRow(context.Background(), queryGetThreadById, id).
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

	if err := r.conn.QueryRow(context.Background(), queryGetThreadBySlug, slug).
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
	if err := r.conn.QueryRow(context.Background(), queryUpdateThreadById, id, req.Title, req.Message).
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
	if err := r.conn.QueryRow(context.Background(), queryUpdateThreadBySlug, slug, req.Title, req.Message).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, threadTime); err != nil {
		return nil, err
	}

	res.Created = strfmt.DateTime(threadTime.UTC()).String()

	return res, nil
}

func (r *ThreadRepository) CreatePosts(forum string, thread int64, posts []*models.RequestNewPost) ([]post_models.ResponsePost, error) {
	query := queryInsertPost

	newPosts := make([]post_models.ResponsePost, 0, len(posts))
	var queryArgs []interface{}
	insertTime := strfmt.DateTime(time.Now())
	for i, post := range posts {
		newPosts[i].Parent = post.Parent
		newPosts[i].Author = post.Author
		newPosts[i].Message = post.Message
		newPosts[i].Forum = forum
		newPosts[i].Thread = thread
		newPosts[i].Created = insertTime.String()

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)

		queryArgs = append(queryArgs, post.Parent, post.Author, post.Message, forum, thread, insertTime)
	}
	query = query[:len(query)-1]
	query += " RETURNING id;"

	rows, err := r.conn.Query(context.Background(), query, queryArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for id := 0; rows.Next(); id++ {
		if err = rows.Scan(&newPosts[id].Id); err != nil {
			return nil, err
		}
	}

	return newPosts, nil
}
func CreateQueryGetPosts(sort string, since string, desc bool, pag *pag_models.Pagination) (string, error) {
	query := queryGetPosts
	orderBy := "ORDER BY created "
	querySince := "AND id > $2"
	limit := pag.Limit

	switch sort {
	case "":
		fallthrough
	case "flat":
		if desc {
			orderBy += "DESC"
		}
		if limit > 0 {
			orderBy += fmt.Sprintf("LIMIT %d", pag.Limit)
		}

		if since != "" {
			query = query + querySince + orderBy

		} else {
			query = query + orderBy
		}
		return query, nil
	case "tree":
	case "parent_tree":
	default:
		return "", thread_repository.SortArgError

	}
	return query, nil
}
func (r *ThreadRepository) GetPostsById(threadId int64, sort string, since string, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error) {
	//orderBy := "ORDER BY created "
	//querySince := "AND id > $2"
	//query := queryGetPosts
	//limit := pag.Limit

	var rows pgx.Rows
	var err error
	query, _ := CreateQueryGetPosts(sort, since, desc, pag)
	//if desc {
	//	orderBy += "DESC"
	//}
	//
	//if limit > 0 {
	//	orderBy += fmt.Sprintf("LIMIT %d", pag.Limit)
	//}

	if since != "" {
		//query = query + querySince + orderBy
		rows, err = r.conn.Query(context.Background(), query, threadId, since)
	} else {
		//query = query + orderBy
		rows, err = r.conn.Query(context.Background(), query, threadId)
	}

	defer rows.Close()

	posts := make([]post_models.ResponsePost, 0)
	for rows.Next() {
		postTime := &time.Time{}
		post := &post_models.ResponsePost{}
		if err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, postTime); err != nil {
			return nil, err
		}
		post.Created = strfmt.DateTime(postTime.UTC()).String()

		posts = append(posts, *post)
	}
	return posts, nil

}
