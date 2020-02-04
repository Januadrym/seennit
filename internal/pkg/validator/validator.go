package validator

import (
	"sync"

	validate "github.com/go-playground/validator"
)

var (
	once      sync.Once
	validator *validate.Validate
)

func New() *validate.Validate {
	once.Do(func() {
		validator = validate.New()
	})
	return validator
}

func Validate(v interface{}) error {
	return New().Struct(v)
}
