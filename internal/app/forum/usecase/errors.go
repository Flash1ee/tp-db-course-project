package forum_usecase

import "github.com/pkg/errors"

var (
	AlreadyExists   = errors.New("forum already exists")
	ForumNotFound   = errors.New("forum not found")
	InternalError   = errors.New("internal error")
	ConstraintError = errors.New("field slug already exists")
)
