package post_postgresql

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
	models4 "tp-db-project/internal/app/forum/models"
	"tp-db-project/internal/app/post/models"
	post_repository "tp-db-project/internal/app/post/repository"
	models3 "tp-db-project/internal/app/thread/models"
	models2 "tp-db-project/internal/app/users/models"
)

const (
	queryCheckPostParent = "SELECT thread from post where id = $1;"
	queryUpdatePost      = "UPDATE post SET message = $2, is_edited = true WHERE id = $1 " +
		"RETURNING id, parent, author, message, is_edited, forum, thread, created;"
	queryUpdatePostNoEdit = "UPDATE post SET is_edited = false WHERE id = $1 " +
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
	conn *pgxpool.Pool
}

func NewPostRepository(conn *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		conn: conn,
	}
}

func (r *PostRepository) Get(id int64, related string) (*models.ResponsePostDetail, error) {
	var err error
	res := &models.ResponsePostDetail{}

	postTime := &time.Time{}
	for _, arg := range strings.Split(related, ",") {
		switch arg {
		case "":
			post := &models.ResponsePost{}
			if err = r.conn.QueryRow(context.Background(), queryGetPost, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, postTime); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}

				return nil, err

			}

			//fmt.Println("POST TIME", postTime)

			//post.Created = strfmt.DateTime(postTime.UTC()).String()
			post.Created = postTime.UTC()
			//fmt.Println("POST TIME", post.Created)

			res.Post = post
		case "user":
			post := &models.ResponsePost{}
			author := &models2.ResponseUser{}
			if err = r.conn.QueryRow(context.Background(), queryGetPostAuthor, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, postTime,
					&author.Nickname, &author.FullName, &author.About, &author.Email); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}
				return nil, err
			}
			//post.Created = postTime.UTC().String()

			post.Created = postTime.UTC()
			res.Post = post
			res.Author = author

		case "thread":
			post := &models.ResponsePost{}
			thread := &models3.ResponseThread{}

			threadTime := &time.Time{}
			if err = r.conn.QueryRow(context.Background(), queryGetPostThread, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, postTime,
					&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
					&thread.Message, &thread.Votes, &thread.Slug, threadTime); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}
				return nil, err
			}
			//post.Created = strfmt.DateTime(postTime.UTC())
			post.Created = postTime.UTC()

			//post.Created = strfmt.DateTime(postTime.UTC()).String()
			thread.Created = strfmt.DateTime(threadTime.UTC()).String()

			res.Post = post
			res.Thread = thread

		case "forum":
			post := &models.ResponsePost{}
			forum := &models4.ResponseForum{}
			if err = r.conn.QueryRow(context.Background(), queryGetPostForum, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, postTime,
					&forum.Title, &forum.User, &forum.Slug, &forum.Posts,
					&forum.Threads); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}
				return nil, err
			}

			//post.Created = strfmt.DateTime(postTime.UTC())
			//
			post.Created = postTime.UTC()

			//post.Created = strfmt.DateTime(postTime.UTC()).String()

			res.Post = post
			res.Forum = forum
		default:
			return nil, post_repository.DefaultErrDB
		}
	}

	return res, nil
}

func (r *PostRepository) Update(id int64, req *models.RequestUpdateMessage) (*models.ResponsePost, error) {
	post := &models.ResponsePost{}
	postTime := &time.Time{}

	if req.Message == "" {
		if err := r.conn.QueryRow(context.Background(), queryUpdatePostNoEdit, id).
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
				&post.IsEdited, &post.Forum, &post.Thread, postTime); err != nil {
			if err == pgx.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err
		}
	} else {
		if err := r.conn.QueryRow(context.Background(), queryUpdatePost, id, req.Message).
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
				&post.IsEdited, &post.Forum, &post.Thread, postTime); err != nil {
			if err == pgx.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err
		}
	}

	post.Created = postTime.UTC()

	return post, nil
}

func (r *PostRepository) CheckParentPost(parent int) (int, error) {
	var threadID int

	err := r.conn.QueryRow(context.Background(), queryCheckPostParent, parent).Scan(&threadID)
	return threadID, err
}
