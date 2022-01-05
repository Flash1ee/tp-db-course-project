package thread_usecase

import (
	"github.com/jackc/pgx/v4"
	"strconv"
	"tp-db-project/internal/app"
	pag_models "tp-db-project/internal/app/models"
	post_models "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/app/thread"
	"tp-db-project/internal/app/thread/models"
	"tp-db-project/internal/app/vote"
	models2 "tp-db-project/internal/app/vote/models"
)

type ThreadUsecase struct {
	repo     thread.Repository
	repoVote vote.Repository
}

func NewThreadUsecase(repo thread.Repository, voteRepo vote.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		repo:     repo,
		repoVote: voteRepo,
	}
}
func (u *ThreadUsecase) GetThreadInfo(slugOrID string) (*models.ResponseThread, error) {
	ID, err := strconv.Atoi(slugOrID)

	switch err {
	case nil:
		res, err := u.repo.GetByID(int64(ID))
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, NotFound
			} else {
				return nil, app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	default:
		res, err := u.repo.GetBySlug(slugOrID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, NotFound
			} else {
				return nil, app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	}
}
func (u *ThreadUsecase) UpdateThread(slugOrID string, req *models.RequestUpdateThread) (*models.ResponseThread, error) {
	ID, err := strconv.Atoi(slugOrID)
	switch err {
	case nil:
		res, err := u.repo.UpdateByID(int64(ID), req)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, NotFound
			} else {
				return nil, app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	default:
		res, err := u.repo.UpdateBySlug(slugOrID, req)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, NotFound
			} else {
				return nil, app.GeneralError{
					Err:         InternalError,
					ExternalErr: err,
				}
			}
		}
		return res, nil
	}
}
func (u *ThreadUsecase) UpdateVoice(slugOrID string, req *models2.RequestVoteUpdate) (bool, error) {
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
		if err == pgx.ErrNoRows {
			return false, NotFound
		} else {
			return false, app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
	}
	isExists, err := u.repoVote.Exists(req.Nickname, th.Id)
	if isExists {
		res, err := u.repoVote.Update(th.Id, req)
		if err != nil {
			return false, app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
		return res, nil
	} else {
		v := &models2.Vote{
			Voice:    req.Voice,
			Nickname: req.Nickname,
			ThreadID: th.Id,
		}
		err = u.repoVote.Create(v)
		if err != nil {
			return false, app.GeneralError{
				Err:         InternalError,
				ExternalErr: err,
			}
		}
		return true, nil
	}
}
func (u *ThreadUsecase) GetPostsBySort(slugOrId string, sort string, since int64, desc bool, pag *pag_models.Pagination) ([]post_models.ResponsePost, error) {
	var idInt int
	var err error
	idInt, err = strconv.Atoi(slugOrId)
	if err != nil {
		th, err := u.repo.GetBySlug(slugOrId)
		if err != nil {
			return nil, err
		} else {
			idInt = int(th.Id)
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
