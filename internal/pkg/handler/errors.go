package handler

import "github.com/pkg/errors"

var (
	InvalidBody       = errors.New("invalid body")
	InvalidParameters = errors.New("invalid parameters in query")
	BDError           = errors.New("can not do bd operation")
)
