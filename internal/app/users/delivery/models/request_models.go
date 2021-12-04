package models

import (
	"github.com/pkg/errors"
)

func validIntFunc(value interface{}) error {
	res, ok := value.(int64)
	if !ok || res < 0 {
		return errors.New("invalid field")
	}
	return nil
}
func validFloatFunc(value interface{}) error {
	res, ok := value.(float64)
	if !ok || res < 0 {
		return errors.New("invalid field")
	}
	return nil
}
