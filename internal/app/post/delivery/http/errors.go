package post_handler

import "github.com/pkg/errors"

var (
	InvalidArgument     = errors.New("invalid argument")
	InvalidBody         = errors.New("invalid body")
	InvalidParamRelated = errors.New("available value of related: user, thread, forum or empty")
	ForumNotFound       = errors.New("forum not found")
)
