// pkg/utils/validator.go

package utils

import "github.com/go-playground/validator/v10"

type Validator interface {
	ValidateStruct(input interface{}) error
}

type val struct {
	*validator.Validate
}

func NewValidator() Validator {
	return &val{validator.New()}
}

func (v *val) ValidateStruct(input interface{}) error {
	return v.Struct(input)
}
