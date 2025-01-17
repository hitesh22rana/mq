// pkg/utils/validator_test.go

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

func TestValidator(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name  string
		fn    func() error
		isErr bool
	}{
		{
			name: "Valid struct",
			fn: func() error {
				return v.ValidateStruct(struct {
					Name   string    `validate:"required"`
					Offset pb.Offset `validate:"required,offset"`
				}{
					Name:   "test",
					Offset: pb.Offset_OFFSET_BEGINNING,
				})
			},
			isErr: false,
		},
		{
			name: "InvalidOffset",
			fn: func() error {
				return v.ValidateStruct(struct {
					Name   string    `validate:"required"`
					Offset pb.Offset `validate:"required,offset"`
				}{
					Name:   "test",
					Offset: pb.Offset_OFFSET_UNKNOWN,
				})
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
