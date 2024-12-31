// pkg/storage/storage.go
//go:generate mockgen -destination=../mocks/mock_storage.go -package=mocks . Storage

package storage

import (
	"errors"
)

const (
	// OffsetBeginning is the offset to start reading messages from the beginning
	OffsetBeginning uint64 = 0

	// OffsetLatest is the offset to start reading messages from the latest
	OffsetLatest uint64 = ^uint64(0)
)

var (
	// ErrInvalidOffset is returned when an invalid offset is provided
	ErrInvalidOffset = errors.New("error: invalid offset provided for message retrieval")

	// ErrInternal is returned when storage is unavailable
	ErrInternal = errors.New("error: storage unavailable")
)

// Storage defines the interface for message storage mechanisms
type Storage interface {
	SaveMessage(string, interface{}) (uint64, error)
	GetMessages(string, string, uint64) ([]interface{}, uint64, error)
	CreateChannel(string) error
	ChannelExists(string) bool
	RemoveChannelFromSubscriberMap(string, string)
}
