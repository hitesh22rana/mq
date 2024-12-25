// pkg/utils/validator_test.go

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name string
		fn   func() error
		err  error
	}{
		{
			name: "ValidateStruct",
			fn: func() error {
				return v.ValidateStruct(struct {
					Name   string `validate:"required"`
					Offset int    `validate:"oneof=0 1"`
				}{
					Name:   "test",
					Offset: 0,
				})
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			assert.Equal(t, tt.err, err)
		})
	}
}
