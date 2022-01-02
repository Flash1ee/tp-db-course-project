package vote_repository

import "github.com/pkg/errors"

var (
	NotFound        = errors.New("Not found")
	InvalidArgument = errors.New("Invalid argument")
)
