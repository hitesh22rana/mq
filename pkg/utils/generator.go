// pkg/utils/generator.go
//go:generate mockgen -destination=../mocks/mock_generator.go -package=mocks . Generator

package utils

import (
	"time"

	"github.com/rs/xid"
)

var (
	messageIDPrefix    = "msg"
	subscriberIDPrefix = "sub"
)

// Generator defines the interface for generating unique IDs and timestamps
type Generator interface {
	GetUniqueMessageID() string
	GetUniqueSubscriberID() string
	GetCurrentTimestamp() int64
}

// generator is the default implementation of the Generator interface
type generator struct{}

// NewGenerator returns a new generator
func NewGenerator() Generator {
	return &generator{}
}

// GetUniqueMessageID returns a unique message ID
func (g *generator) GetUniqueMessageID() string {
	return messageIDPrefix + xid.New().String()
}

// GetUniqueSubscriberID returns a unique subscriber ID
func (g *generator) GetUniqueSubscriberID() string {
	return subscriberIDPrefix + xid.New().String()
}

// GetCurrentTimestamp returns the current timestamp
func (g *generator) GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
