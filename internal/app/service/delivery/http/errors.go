package service_handler

import "github.com/pkg/errors"

var (
	NotFound    = errors.New("user not found")
	ConflictErr = errors.New("conflict update data - already exists field")
)
