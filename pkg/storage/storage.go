// pkg/storage/storage.go

package storage

import (
	"errors"
)

var (
	// ErrInvalidOffset is returned when an invalid offset is provided
	ErrInvalidOffset = errors.New("error: invalid offset provided for message retrieval")

	// ErrInternal is returned when storage is unavailable
	ErrInternal = errors.New("error: storage unavailable")
)

// Storage defines the interface for message storage mechanisms
type Storage interface {
	SaveMessage(string, interface{}) (int64, error)
	GetMessages(string, int64) ([]interface{}, int64, error)
	CreateChannel(string) error
	ChannelExists(string) bool
	DeleteChannel(string) error
}
