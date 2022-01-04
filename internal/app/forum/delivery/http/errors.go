package forum_handler

import "github.com/pkg/errors"

var (
	InvalidArgument = errors.New("invalid argument")
	InvalidBody     = errors.New("invalid body")

	ForumNotFound        = errors.New("forum not found")
	InvalidLimitArgument = errors.New("invalid limit argument")
	InvalidDescArgument  = errors.New("invalid desc argument")
)
