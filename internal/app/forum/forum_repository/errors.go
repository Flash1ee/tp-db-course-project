package forum_repository

import (
	"github.com/pkg/errors"
	"tp-db-project/internal/app"
)

var (
	DefaultErrDB = errors.New("something wrong DB")
	ArgError     = errors.New("invalid argument")
)

func NewDBError(externalErr error) *app.GeneralError {
	return &app.GeneralError{
		Err:         DefaultErrDB,
		ExternalErr: externalErr,
	}

}
