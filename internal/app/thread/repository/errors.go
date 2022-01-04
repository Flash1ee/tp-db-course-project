package repository

import "github.com/pkg/errors"

var (
	ArgError     = errors.New("invalid argument")
	NotFound     = errors.New("Not found")
	SortArgError = errors.New("invalid sort flag")
)
