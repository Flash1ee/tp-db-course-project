package post_postgresql

import (
	"context"
	"database/sql"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"time"
	"tp-db-project/internal/app/post/models"
	post_repository "tp-db-project/internal/app/post/repository"
)

const (
	queryUpdatePost = "UPDATE post SET message = $2, is_edited = true WHERE id = $1 " +
		"RETURNING id, parent, author, message, is_edited, forum, thread, created;"
	queryGetPost       = "SELECT id, parent, author, message, is_edited, forum, thread, created from post WHERE id = $1;"
	queryGetPostAuthor = "SELECT post.id, post.parent, post.author, post.message, post.is_edited, post.forum, post.thread, post.created, " +
		"a.nickname, a.fullname, a.about, a.email from post " +
		"JOIN users a on a.nickname = post.author WHERE post.id = $1;"
	queryGetPostThread = "SELECT post.id, post.parent, post.author, post.message, post.is_edited, post.forum, post.thread, post.created, " +
		"th.id, th.title, th.author, th.forum, th.message, th.votes, th.slug, th.created from post " +
		"JOIN thread th on th.id = post.thread WHERE post.id = $1;"
	queryGetPostForum = "SELECT post.id, post.parent, post.author, post.message, post.is_edited, post.forum, post.thread, post.created, " +
		"f.title, f.users_nickname, f.slug, f.posts, f.threads from post " +
		"JOIN forum f on f.slug = post.forum WHERE post.id = $1;"
)

type PostRepository struct {
	conn *pgx.Conn
}

func NewPostRepository(conn *pgx.Conn) *PostRepository {
	return &PostRepository{
		conn: conn,
	}
}

func (r *PostRepository) Get(id int64, related string) (*models.ResponsePostDetail, error) {
	var err error
	res := &models.ResponsePostDetail{}

	postTime := &time.Time{}
	switch related {
	case "":
		if err = r.conn.QueryRow(context.Background(), queryGetPost, id).
			Scan(&res.Post.Id, &res.Post.Parent, &res.Post.Author, &res.Post.Message,
				&res.Post.IsEdited, &res.Post.Forum, &res.Post.Thread, postTime); err != nil {
			if err == sql.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err

		}
		res.Post.Created = strfmt.DateTime(postTime.UTC()).String()

	case "user":
		if err = r.conn.QueryRow(context.Background(), queryGetPostAuthor, id).
			Scan(&res.Post.Id, &res.Post.Parent, &res.Post.Author, &res.Post.Message,
				&res.Post.IsEdited, &res.Post.Forum, &res.Post.Thread, postTime,
				&res.Author.Nickname, &res.Author.FullName, &res.Author.About, &res.Author.Email); err != nil {
			if err == sql.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err

		}
		res.Post.Created = strfmt.DateTime(postTime.UTC()).String()

	case "thread":
		threadTime := &time.Time{}
		if err = r.conn.QueryRow(context.Background(), queryGetPostThread, id).
			Scan(&res.Post.Id, &res.Post.Parent, &res.Post.Author, &res.Post.Message,
				&res.Post.IsEdited, &res.Post.Forum, &res.Post.Thread, postTime,
				&res.Thread.Id, &res.Thread.Title, &res.Thread.Author, &res.Thread.Forum,
				&res.Thread.Message, &res.Thread.Votes, &res.Thread.Slug, threadTime); err != nil {
			if err == sql.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err
		}

		res.Post.Created = strfmt.DateTime(postTime.UTC()).String()
		res.Thread.Created = strfmt.DateTime(threadTime.UTC()).String()

	case "forum":
		if err = r.conn.QueryRow(context.Background(), queryGetPostForum, id).
			Scan(&res.Post.Id, &res.Post.Parent, &res.Post.Author, &res.Post.Message,
				&res.Post.IsEdited, &res.Post.Forum, &res.Post.Thread, postTime,
				&res.Forum.Title, &res.Forum.User, &res.Forum.Slug, &res.Forum.Posts,
				&res.Forum.Threads); err != nil {
			if err == sql.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err
		}

		res.Post.Created = strfmt.DateTime(postTime.UTC()).String()
	default:
		return nil, post_repository.DefaultErrDB
	}

	return res, nil
}

func (r *PostRepository) Update(id int64, req *models.RequestUpdateMessage) (*models.ResponsePost, error) {
	var post *models.ResponsePost

	postTime := &time.Time{}
	if err := r.conn.QueryRow(context.Background(), queryUpdatePost, id, req.Message).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, postTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, post_repository.NotFound
		}
		return nil, err
	}
	post.Created = strfmt.DateTime(postTime.UTC()).String()

	return post, nil
}
