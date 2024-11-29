// pkg/storage/memory_storage.go

package storage

import (
	"fmt"
	"sync"
	
	"go.uber.org/zap"
)

// MemoryStorage is an in-memory implementation of the Storage interface
type MemoryStorage struct {
	mu     sync.RWMutex
	logger *zap.Logger
	data   map[string][]interface{}
}

// NewMemoryStorage initializes a new MemoryStorage instance
func NewMemoryStorage(logger *zap.Logger) *MemoryStorage {
	return &MemoryStorage{
		mu:     sync.RWMutex{},
		logger: logger,
		data:   make(map[string][]interface{}),
	}
}

// SaveMessage saves a message to the specified channel
func (m *MemoryStorage) SaveMessage(channel string, message interface{}) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the channel exists
	if _, exists := m.data[channel]; !exists {
		m.logger.Info("channel does not exist", zap.String("channel", channel))

		// Create a new channel
		if err := m.CreateChannel(channel); err != nil {
			return 0, ErrInternal
		}
	}

	// Append the message to the channel
	m.data[channel] = append(m.data[channel], message)

	// Return the index of the message in the channel
	return int64(len(m.data[channel]) - 1), nil
}

// GetMessages retrieves all messages from the specified channel
func (m *MemoryStorage) GetMessages(channel string, offset int64) ([]interface{}, error) {
	if offset < 0 {
		return nil, ErrInvalidOffset
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, exists := m.data[channel]
	if !exists {
		m.logger.Info("channel does not exist", zap.String("channel", channel))
		return nil, fmt.Errorf("channel '%s' does not exist", channel)
	}

	// Check if the offset is greater than the number of messages
	if int(offset) >= len(messages) {
		return nil, ErrInvalidOffset
	}

	// Return the copy of the messages to prevent data mutation
	return append([]interface{}(nil), messages[offset:]...), nil
}

// CreateChannel creates a new channel
func (m *MemoryStorage) CreateChannel(channel string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[channel]; exists {
		m.logger.Info("channel does not exist", zap.String("channel", channel))
		return fmt.Errorf("channel '%s' already exists", channel)
	}

	m.data[channel] = []interface{}{}
	return nil
}

// DeleteChannel deletes an existing channel
func (m *MemoryStorage) DeleteChannel(channel string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[channel]; !exists {
		m.logger.Info("channel does not exist", zap.String("channel", channel))
		return fmt.Errorf("channel '%s' does not exist", channel)
	}

	delete(m.data, channel)
	return nil
}
