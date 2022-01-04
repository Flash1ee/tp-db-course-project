package thread_usecase

import "github.com/pkg/errors"

var (
	InternalError = errors.New("internal error")
	NotFound      = errors.New("thread not found")
)
