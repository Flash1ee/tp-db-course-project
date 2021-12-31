package post_repository

import (
	"github.com/pkg/errors"
	"tp-db-project/internal/app"
)

var (
	DefaultErrDB = errors.New("something wrong DB")
	NotFound     = errors.New("Not found")
)

func NewDBError(externalErr error) *app.GeneralError {
	return &app.GeneralError{
		Err:         DefaultErrDB,
		ExternalErr: externalErr,
	}

}
