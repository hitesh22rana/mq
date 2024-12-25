// pkg/utils/generator_test.go

package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	g := NewGenerator()

	tests := []struct {
		name string
		fn   func() string
	}{
		{
			name: "GetUniqueMessageID",
			fn:   g.GetUniqueMessageID,
		},
		{
			name: "GetUniqueSubscriberID",
			fn:   g.GetUniqueSubscriberID,
		},
		{
			name: "GetCurrentTimestamp",
			fn: func() string {
				timeStamp := g.GetCurrentTimestamp()
				return strconv.FormatInt(timeStamp, 10)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.fn()
			assert.NotEmpty(t, res)
		})
	}
}
