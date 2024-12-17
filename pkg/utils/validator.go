// pkg/utils/validator.go

package utils

import "github.com/go-playground/validator/v10"

// Validator defines the interface for validating input structs
type Validator interface {
	ValidateStruct(input interface{}) error
}

// val is the default implementation of the Validator interface
type val struct {
	*validator.Validate
}

// NewValidator returns a new validator
func NewValidator() Validator {
	return &val{validator.New()}
}

// ValidateStruct validates the input struct
func (v *val) ValidateStruct(input interface{}) error {
	return v.Struct(input)
}
