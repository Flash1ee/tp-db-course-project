package users_handler

import "github.com/pkg/errors"

var (
	InvalidBody     = errors.New("invalid body in request")
	InvalidArgument = errors.New("invalid argument in request url")
	NotFound        = errors.New("user not found")
	ConflictErr     = errors.New("conflict update data - already exists field")
)
