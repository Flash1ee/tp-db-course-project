package post_postgresql

import (
	"github.com/jackc/pgx"
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
	conn *pgx.ConnPool
}

func NewPostRepository(conn *pgx.ConnPool) *PostRepository {
	return &PostRepository{
		conn: conn,
	}
}

func (r *PostRepository) Get(id int64, related string) (*models.ResponsePostDetail, error) {
	var err error
	res := &models.ResponsePostDetail{}

	postTime := time.Time{}
	for _, arg := range strings.Split(related, ",") {
		switch arg {
		case "":
			post := &models.ResponsePost{}
			if err = r.conn.QueryRow(queryGetPost, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, &postTime); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}

				return nil, err

			}

			//fmt.Println("POST TIME", postTime)
			//t := strfmt.DateTime(postTime.UTC())
			//post.Created = &t
			//post.Created = time.Time(*postTime).String()
			post.Created = postTime.Format(time.RFC3339)
			//fmt.Println("POST TIME", post.Created)
			//fmt.Println(post.Created)
			//post.Created = post.Created.Format(time.RFC3339)
			res.Post = post
		case "user":
			post := &models.ResponsePost{}
			author := &models2.ResponseUser{}
			if err = r.conn.QueryRow(queryGetPostAuthor, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, &postTime,
					&author.Nickname, &author.FullName, &author.About, &author.Email); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}
				return nil, err
			}
			post.Created = postTime.Format(time.RFC3339)

			//post.Created = postTime.UTC()
			//post.Created = time.Time(*postTime).String()

			//post.Created = time.Time(*postTime).UTC().Truncate(time.Nanosecond)
			//post.Created = *postTime
			//t := strfmt.DateTime(postTime.UTC())
			//post.Created = &t
			res.Post = post
			res.Author = author

		case "thread":
			post := &models.ResponsePost{}
			thread := &models3.ResponseThread{}

			//threadTime := &strfmt.DateTime{}
			if err = r.conn.QueryRow(queryGetPostThread, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, &postTime,
					&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
					&thread.Message, &thread.Votes, &thread.Slug, &thread.Created); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}
				return nil, err
			}
			//post.Created = strfmt.DateTime(postTime.UTC())
			//post.Created = time.Time(*postTime).String()

			//post.Created = time.Time(*postTime).UTC().Truncate(time.Nanosecond)

			post.Created = postTime.Format(time.RFC3339)
			//t := strfmt.DateTime(postTime.UTC())
			//post.Created = &t
			//post.Created = *postTime

			//post.Created = strfmt.DateTime(postTime.UTC()).String()
			//thread.Created = *threadTime.String()
			//thread.Created = time.Time(*threadTime).UTC().Truncate(time.Nanosecond)

			res.Post = post
			res.Thread = thread

		case "forum":
			post := &models.ResponsePost{}
			forum := &models4.ResponseForum{}
			if err = r.conn.QueryRow(queryGetPostForum, id).
				Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
					&post.IsEdited, &post.Forum, &post.Thread, &postTime,
					&forum.Title, &forum.User, &forum.Slug, &forum.Posts,
					&forum.Threads); err != nil {
				if err == pgx.ErrNoRows {
					return nil, post_repository.NotFound
				}
				return nil, err
			}

			//post.Created = strfmt.DateTime(postTime.UTC())
			//
			//post.Created = postTime.UTC()
			//t := strfmt.DateTime(postTime.UTC())
			//post.Created = &t
			//post.Created = *postTime

			//post.Created = strfmt.DateTime(postTime.UTC()).String()
			post.Created = postTime.Format(time.RFC3339)
			//post.Created = time.Time(*postTime).String()

			//post.Created = time.Time(*postTime).UTC().Truncate(time.Nanosecond)
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
	postTime := time.Time{}
	//postTime := &strfmt.DateTime{}

	if req.Message == "" {
		if err := r.conn.QueryRow(queryUpdatePostNoEdit, id).
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
				&post.IsEdited, &post.Forum, &post.Thread, &postTime); err != nil {
			if err == pgx.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err
		}
	} else {
		if err := r.conn.QueryRow(queryUpdatePost, id, req.Message).
			Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
				&post.IsEdited, &post.Forum, &post.Thread, &postTime); err != nil {
			if err == pgx.ErrNoRows {
				return nil, post_repository.NotFound
			}
			return nil, err
		}
	}

	//post.Created = postTime.UTC()
	//t := strfmt.DateTime(postTime.UTC())
	//post.Created = &t
	post.Created = postTime.Format(time.RFC3339)
	//post.Created = time.Time(*postTime).String()

	//post.Created = time.Time(*postTime).UTC().Truncate(time.Nanosecond)

	return post, nil
}

func (r *PostRepository) CheckParentPost(parent int) (int, error) {
	var threadID int

	err := r.conn.QueryRow(queryCheckPostParent, parent).Scan(&threadID)
	return threadID, err
}
