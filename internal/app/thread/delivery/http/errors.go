package thread_handler

import "github.com/pkg/errors"

var (
	InvalidArgument = errors.New("invalid argument")
	InvalidBody     = errors.New("invalid body")
	InvalidVoice    = errors.New("available voice is 1 or -1")
	ForumNotFound   = errors.New("forum not found")
)
