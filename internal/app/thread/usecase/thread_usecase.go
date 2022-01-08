package thread_usecase

import (
	"github.com/jackc/pgx"
	"strconv"
	"tp-db-project/internal/app"
	pag_models "tp-db-project/internal/app/models"
	"tp-db-project/internal/app/post"
	post_models "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/app/thread"
	"tp-db-project/internal/app/thread/models"
	"tp-db-project/internal/app/thread/repository"
	"tp-db-project/internal/app/users"
	"tp-db-project/internal/app/vote"
	models2 "tp-db-project/internal/app/vote/models"
)

type ThreadUsecase struct {
	repo     thread.Repository
	repoVote vote.Repository
	repoPost post.Repository
	repoUser users.Repository
}

func NewThreadUsecase(repo thread.Repository, voteRepo vote.Repository, postRepo post.Repository,
	userRepo users.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		repo:     repo,
		repoVote: voteRepo,
		repoPost: postRepo,
		repoUser: userRepo,
	}
}
func (u *ThreadUsecase) GetThreadInfo(slugOrID string) (*models.ResponseThread, error) {
	ID, err := strconv.Atoi(slugOrID)

	switch err {
	case nil:
		res, err := u.repo.GetByID(int64(ID))
		if err != nil {
			if err == repository.NotFound {

				return nil, &app.GeneralError{
					Err:         NotFound,
					ExternalErr: err,
				}
			} else {
				return nil, &app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	default:
		res, err := u.repo.GetBySlug(slugOrID)
		if err != nil {
			if err == repository.NotFound {

				return nil, &app.GeneralError{
					Err:         NotFound,
					ExternalErr: err,
				}
			}
			return nil, &app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
		return res, nil
	}
}
func (u *ThreadUsecase) UpdateThread(slugOrID string, req *models.RequestUpdateThread) (*models.ResponseThread, error) {
	ID, err := strconv.Atoi(slugOrID)
	var th *models.ResponseThread
	switch err {
	case nil:
		th, err = u.repo.GetByID(int64(ID))
		if err != nil {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		}
		if req.Title == "" {
			req.Title = th.Title
		}
		if req.Message == "" {
			req.Message = th.Message
		}
		res, err := u.repo.UpdateByID(int64(ID), req)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, NotFound
			} else {
				return nil, &app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	default:
		th, err = u.repo.GetBySlug(slugOrID)
		if err != nil {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		}
		if req.Title == "" {
			req.Title = th.Title
		}
		if req.Message == "" {
			req.Message = th.Message
		}
		res, err := u.repo.UpdateBySlug(slugOrID, req)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, NotFound
			} else {
				return nil, &app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	}
}
func (u *ThreadUsecase) UpdateVoice(slugOrID string, req *models2.RequestVoteUpdate) (*models.ResponseThread, error) {
	var err error
	var th *models.ResponseThread
	ID, err := strconv.Atoi(slugOrID)

	switch err {
	case nil:
		th, err = u.repo.GetByID(int64(ID))
	default:
		th, err = u.repo.GetBySlug(slugOrID)
	}

	if err != nil {
		if err == repository.NotFound {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		} else {
			return nil, &app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
	}
	if _, err := u.repoUser.GetByNickname(req.Nickname); err != nil {
		return nil, &app.GeneralError{
			Err:         AuthorNotFound,
			ExternalErr: err,
		}
	}
	isExists, err := u.repoVote.Exists(req.Nickname, th.Id)
	if isExists {
		ok, err := u.repoVote.Update(th.Id, req)
		if err != nil {
			return nil, &app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
		if ok {
			th.Votes += req.Voice * 2
		}
		return th, nil
	} else {
		v := &models2.Vote{
			Voice:    req.Voice,
			Nickname: req.Nickname,
			ThreadID: th.Id,
		}
		err = u.repoVote.Create(v)
		if err != nil {
			return nil, &app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
		th.Votes += req.Voice
		return th, nil
	}
}
func (u *ThreadUsecase) GetPostsBySort(slugOrId string, sort string, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error) {
	var idInt int
	var err error
	idInt, err = strconv.Atoi(slugOrId)
	if err != nil {
		th, err := u.repo.GetBySlug(slugOrId)
		if err != nil {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		} else {
			idInt = int(th.Id)
		}
	} else {
		if _, err := u.repo.GetByID(int64(idInt)); err != nil {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		}
	}

	switch sort {
	case "flat":
		return u.repo.GetPostsByFlats(idInt, since, desc, pag)
	case "tree":
		return u.repo.GetPostsByTree(idInt, since, desc, pag)
	case "parent_tree":
		return u.repo.GetPostsByParentTree(idInt, since, desc, pag)
	default:
		return u.repo.GetPostsByFlats(idInt, since, desc, pag)
	}
}
func (u *ThreadUsecase) CreatePosts(slugOrID string, posts []*models.RequestNewPost) ([]post_models.ResponsePost, error) {
	var idInt int
	var forumName string
	var err error
	var th *models.ResponseThread
	idInt, err = strconv.Atoi(slugOrID)
	if err != nil {
		th, err = u.repo.GetBySlug(slugOrID)
		if err != nil {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		} else {
			idInt = int(th.Id)
		}
	} else {
		th, err = u.repo.GetByID(int64(idInt))
		if err != nil {
			return nil, &app.GeneralError{
				Err:         NotFound,
				ExternalErr: err,
			}
		}
	}
	forumName = th.Forum
	if len(posts) == 0 {
		return nil, &app.GeneralError{
			Err:         CreatedEmpty,
			ExternalErr: err,
		}
	}
	if posts[0].Parent != 0 {
		parentTh, err := u.repoPost.CheckParentPost(int(posts[0].Parent))
		if err != nil {
			return nil, &app.GeneralError{
				Err:         ParentNotFound,
				ExternalErr: err,
			}
		}

		if parentTh != idInt {
			return nil, &app.GeneralError{
				Err:         ParentInvalid,
				ExternalErr: err,
			}
		}
	}
	_, err = u.repoUser.GetByNickname(posts[0].Author)
	if err != nil {
		return nil, &app.GeneralError{
			Err:         AuthorNotFound,
			ExternalErr: err,
		}
	}
	res, err := u.repo.CreatePosts(forumName, int64(idInt), posts)
	if err != nil {
		return nil, &app.GeneralError{
			Err:         err,
			ExternalErr: err,
		}
	}
	return res, nil
}
