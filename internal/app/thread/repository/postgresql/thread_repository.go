package thread_postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	models2 "tp-db-project/internal/app/forum/models"
	pag_models "tp-db-project/internal/app/models"
	post_models "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/app/thread/models"
	"tp-db-project/internal/app/thread/repository"
)

const (
	queryGetThreadById = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE id = $1;"
	queryGetThreadBySlug = "SELECT id, title, author, forum, message, votes, slug, created FROM thread " +
		"WHERE slug = $1;"
	queryCreateForumThreadWithTime = "INSERT INTO thread(title, author, forum, message, slug, created) VALUES($1, $2, $3, $4, $5, $6) " +
		"RETURNING id, title, author, forum, message, votes, slug, created;"
	queryCreateForumThreadNoTime = "INSERT INTO thread(title, author, forum, message, slug) VALUES($1, $2, $3, $4, $5) " +
		"RETURNING id, title, author, forum, message, votes, slug, created; "
	queryUpdateThreadById = "UPDATE thread SET title = $2, message = $3 " +
		"WHERE id = $1 RETURNING id, title, author, forum, message, votes, slug, created;"
	queryUpdateThreadBySlug = "UPDATE thread SET title = $2, message = $3 " +
		"WHERE slug = $1 RETURNING id, title, author, forum, message, votes, slug, created;"

	queryInsertPost = "INSERT INTO post(parent, author, message, forum, thread, created) VALUES "
)

type ThreadRepository struct {
	conn *pgxpool.Pool
}

func NewThreadRepository(conn *pgxpool.Pool) *ThreadRepository {
	return &ThreadRepository{
		conn: conn,
	}
}

func (r *ThreadRepository) CreateThread(req *models2.RequestCreateThread) (*models.ResponseThread, error) {
	var err error
	res := &models.ResponseThread{}
	if req == nil {
		return nil, repository.ArgError
	}
	//threadTime := &time.Time{}
	//threadTime := &strfmt.DateTime{}
	if req.Created == "" {
		err = r.conn.QueryRow(context.Background(), queryCreateForumThreadNoTime, req.Title, req.Author, req.Forum, req.Message, req.Slug).
			Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
				&res.Message, &res.Votes, &res.Slug, &res.Created)
	} else {
		err = r.conn.QueryRow(context.Background(), queryCreateForumThreadWithTime, req.Title, req.Author, req.Forum, req.Message, req.Slug, req.Created).
			Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
				&res.Message, &res.Votes, &res.Slug, &res.Created)
	}

	//fmt.Printf("THREAD SLUG AFTER CREATED = %s FORUM = %s\n", res.Slug, res.Forum)
	//res.Created = strfmt.DateTime(threadTime.UTC()).String()

	//res.Created = time.Time(*threadTime).UTC()
	return res, err
}

func (r *ThreadRepository) GetByID(id int64) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	//threadTime := &time.Time{}
	//threadTime := &strfmt.DateTime{}



	if err := r.conn.QueryRow(context.Background(), queryGetThreadById, id).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, err
	}
	//res.Created = strfmt.DateTime(threadTime.UTC()).String()
	//res.Created = time.Time(*threadTime).UTC()
	return res, nil
}

func (r *ThreadRepository) GetBySlug(slug string) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	//threadTime := &strfmt.DateTime{}

	//threadTime := &time.Time{}

	if err := r.conn.QueryRow(context.Background(), queryGetThreadBySlug, slug).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NotFound
		}

		return nil, err
	}
	//res.Created = strfmt.DateTime(threadTime.UTC()).String()

	//res.Created = time.Time(*threadTime).UTC()
	return res, nil
}

func (r *ThreadRepository) UpdateByID(id int64, req *models.RequestUpdateThread) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	//threadTime := &time.Time{}
	//threadTime := &strfmt.DateTime{}
	if err := r.conn.QueryRow(context.Background(), queryUpdateThreadById, id, req.Title, req.Message).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		return nil, err
	}

	//res.Created = strfmt.DateTime(threadTime.UTC()).String()
	//res.Created = time.Time(*threadTime).UTC()

	return res, nil
}

func (r *ThreadRepository) UpdateBySlug(slug string, req *models.RequestUpdateThread) (*models.ResponseThread, error) {
	res := &models.ResponseThread{}
	//threadTime := &time.Time{}
	//threadTime := &strfmt.DateTime{}

	if err := r.conn.QueryRow(context.Background(), queryUpdateThreadBySlug, slug, req.Title, req.Message).
		Scan(&res.Id, &res.Title, &res.Author, &res.Forum,
			&res.Message, &res.Votes, &res.Slug, &res.Created); err != nil {
		return nil, err
	}

	//res.Created = strfmt.DateTime(threadTime.UTC()).String()
	//res.Created = time.Time(*threadTime).UTC()

	return res, nil
}

func (r *ThreadRepository) CreatePosts(forum string, thread int64, posts []*models.RequestNewPost) ([]post_models.ResponsePost, error) {
	query := queryInsertPost

	newPosts := make([]post_models.ResponsePost, len(posts), len(posts))
	queryArgs := make([]interface{}, 0, 0)
	//insertTime := strfmt.DateTime(time.Now())
	insertTime := time.Now()

	for i, post := range posts {
		newPosts[i].Parent = post.Parent
		newPosts[i].Author = post.Author
		newPosts[i].Message = post.Message
		newPosts[i].Forum = forum
		newPosts[i].Thread = thread
		//newPosts[i].Created = insertTime.String()
		newPosts[i].Created = insertTime.Format(time.RFC3339)

		//newPosts[i].Created = time.Time(insertTime)

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)

		queryArgs = append(queryArgs, post.Parent, post.Author, post.Message, forum, thread, insertTime)
		//queryArgs = append(queryArgs, post.Parent, post.Author, post.Message, forum, thread)

	}
	query = query[:len(query)-1]
	query += " RETURNING id;"

	rows, err := r.conn.Query(context.Background(), query, queryArgs...)
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

func (r *ThreadRepository) GetPostsByFlats(id int, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error) {
	queryGetPostsByFlat := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM post " +
		"WHERE thread = $1 "
	var rows pgx.Rows
	var err error
	query := queryGetPostsByFlat
	if since != -1 && desc {
		query += " and id < $2"
	} else if since != -1 && !desc {
		query += " and id > $2"
	} else if since != -1 {
		query += " and id > $2"
	}
	if desc {
		query += " ORDER BY created desc, id desc"
	} else if !desc {
		query += " ORDER BY created asc, id"
	} else {
		query += " ORDER BY created, id"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", pag.Limit)

	if since == -1 {
		rows, err = r.conn.Query(context.Background(), query, id)
	} else {
		rows, err = r.conn.Query(context.Background(), query, id, since)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]post_models.ResponsePost, 0, 0)
	for rows.Next() {
		post := post_models.ResponsePost{}
		//timeTmp := &strfmt.DateTime{}

		timeTmp := time.Now()

		err = rows.Scan(&post.Id, &post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&timeTmp)
		if err != nil {
			return nil, err
		}

		//t := strfmt.DateTime(timeTmp.UTC())
		//post.Created = &t
		//post.Created = time.Time(*timeTmp).String()

		//post.Created = time.Time(*timeTmp).UTC().Truncate(time.Nanosecond)

		//post.Created = strfmt.DateTime(timeTmp.UTC()).String()
		post.Created = timeTmp.Format(time.RFC3339)

		posts = append(posts, post)
	}

	return posts, err
}

func (r *ThreadRepository) GetPostsByTree(id int, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error) {
	queryGetPostsByTree := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM post " +
		"WHERE thread = $1 "
	//queryGetPostsByTreeSince := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM post " +
	//	"WHERE thread = $1 AND path < (SELECT path FROM post WHERE id = $2) "
	var rows pgx.Rows
	var err error
	//orderBy := "ORDER BY path %s, id %s"
	//qLimit := "LIMIT $2"

	query := queryGetPostsByTree
	if since != -1 && desc {
		query += " and path < "
	} else if since != -1 && !desc {
		query += " and path > "
	} else if since != -1 {
		query += " and path > "
	}
	if since != -1 {
		query += fmt.Sprintf(` (SELECT path FROM post WHERE id = %d) `, since)
	}
	if desc {
		query += " ORDER BY path desc"
	} else if !desc {
		query += " ORDER BY path asc, id"
	} else {
		query += " ORDER BY path, id"
	}
	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", pag.Limit)

	//if desc {
	//	orderBy = fmt.Sprintf(orderBy, "DESC", "DESC")
	//} else {
	//	orderBy = fmt.Sprintf(orderBy, "ASC", "ASC")
	//}
	//if since == -1 {
	rows, err = r.conn.Query(context.Background(), query, id)
	//} else {
	//	rows, err = r.conn.Query(context.Background(), query, id, since)
	//}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]post_models.ResponsePost, 0, 0)
	for rows.Next() {
		post := post_models.ResponsePost{}
		timeTmp := time.Time{}
		//timeTmp := &strfmt.DateTime{}
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&timeTmp)
		if err != nil {
			return nil, err
		}
		//post.Created = strfmt.DateTime(timeTmp.UTC())
		//t := strfmt.DateTime(timeTmp.UTC())
		//post.Created = &t
		post.Created = timeTmp.Format(time.RFC3339)
		//post.Created = time.Time(*timeTmp).UTC().Truncate(time.Nanosecond)
		//post.Created = time.Time(*timeTmp).String()


		//post.Created = strfmt.DateTime(timeTmp.UTC()).String()

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *ThreadRepository) GetPostsByParentTree(id int, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error) {
	//queryGetPostsByParentTree := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM post " +
	//	"WHERE thread = $1 "
	//queryGetPostsByParentTreeSince := "SELECT id, parent, author, message, is_edited, forum, thread, created FROM post " +
	//	"WHERE thread = $1 AND path < (SELECT path FROM post WHERE id = $2) "
	var rows pgx.Rows
	var err error
	//orderBy := "ORDER BY path %s, id %s"
	//qLimit := "LIMIT $2"
	//if desc {
	//	orderBy = fmt.Sprintf(orderBy, "DESC", "DESC")
	//} else {
	//	orderBy = fmt.Sprintf(orderBy, "ASC", "ASC")
	//
	//}

	if since == -1 {
		if desc {
			rows, err = r.conn.Query(context.Background(), `
					SELECT id, parent, author, message, is_edited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, id ASC;`,
				id,
				pag.Limit)
		} else {
			rows, err = r.conn.Query(context.Background(), `
					SELECT id, parent, author, message, is_edited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 ORDER BY id ASC LIMIT $2)
					ORDER BY path ASC, id ASC;`,
				id,
				pag.Limit)
		}
	} else {
		if desc {
			rows, err = r.conn.Query(context.Background(), `
					SELECT id, parent, author, message, is_edited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 AND path[1] <
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY path[1] DESC, path ASC, id ASC;`,
				id,
				since,
				pag.Limit)
		} else {
			rows, err = r.conn.Query(context.Background(), `
					SELECT id, parent, author, message, is_edited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 AND path[1] >
					(SELECT path[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path ASC, id ASC;`,
				id,
				since,
				pag.Limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]post_models.ResponsePost, 0, 0)
	for rows.Next() {
		post := post_models.ResponsePost{}
		timeTmp := time.Time{}
		//timeTmp := &strfmt.DateTime{}


		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&timeTmp)
		if err != nil {
			return nil, err
		}
		//t := strfmt.DateTime(timeTmp.UTC())
		//post.Created = &t
		post.Created = timeTmp.Format(time.RFC3339)
		//post.Created = time.Time(*timeTmp).String()

		//post.Created = time.Time(*timeTmp).UTC().Truncate(time.Nanosecond)
		//post.Created = strfmt.DateTime(timeTmp.UTC()).String()

		posts = append(posts, post)
	}

	return posts, nil
}
