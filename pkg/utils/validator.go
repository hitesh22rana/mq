// pkg/utils/validator.go
//go:generate mockgen -destination=../mocks/mock_validator.go -package=mocks . Validator

package utils

import (
	"github.com/go-playground/validator/v10"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// Validator defines the interface for validating input structs
type Validator interface {
	ValidateStruct(input interface{}) error
}

// val is the default implementation of the Validator interface
type val struct {
	*validator.Validate
}

func validateOffset(fl validator.FieldLevel) bool {
	offset := fl.Field().Interface().(pb.Offset)
	switch offset {
	case pb.Offset_OFFSET_BEGINNING, pb.Offset_OFFSET_LATEST:
		return true
	default:
		return false
	}
}

// NewValidator returns a new validator
func NewValidator() Validator {
	_val := validator.New()

	// Register custom offset validation
	_val.RegisterValidation("offset", validateOffset)
	return &val{_val}
}

// ValidateStruct validates the input struct
func (v *val) ValidateStruct(input interface{}) error {
	return v.Struct(input)
}
