package users_usecase

import "github.com/pkg/errors"

var (
	AlreadyExistsErr = errors.New("user already exists")
	NotFound         = errors.New("user not found")
	ConstraintError  = errors.New("invalid params - email already exists")
)
