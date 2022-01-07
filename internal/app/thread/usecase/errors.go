package thread_usecase

import "github.com/pkg/errors"

var (
	CreatedEmpty   = errors.New("empty query")
	AuthorNotFound = errors.New("author not found")
	ParentInvalid  = errors.New("parent post was created in another thread")
	ParentNotFound = errors.New("parent not found")
	InternalError  = errors.New("internal error")
	NotFound       = errors.New("thread not found")
)
