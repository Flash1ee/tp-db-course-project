package thread_repository

import "github.com/pkg/errors"

var (
	NotFound     = errors.New("Not found")
	SortArgError = errors.New("invalid sort flag")
)
