// pkg/storage/memory_storage.go

package storage

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// MemoryStorageOptions represents the options for the MemoryStorage
type MemoryStorageOptions struct {
	BatchSize int
}

// MemoryStorage is an in-memory implementation of the Storage interface
type MemoryStorage struct {
	mu        sync.RWMutex
	logger    *zap.Logger
	batchSize int
	data      map[string][]interface{}
}

// NewMemoryStorage initializes a new MemoryStorage instance
func NewMemoryStorage(logger *zap.Logger, options *MemoryStorageOptions) *MemoryStorage {
	return &MemoryStorage{
		mu:        sync.RWMutex{},
		logger:    logger,
		batchSize: options.BatchSize,
		data:      make(map[string][]interface{}),
	}
}

// SaveMessage saves a message to the specified channel
func (m *MemoryStorage) SaveMessage(channel string, message interface{}) (int64, error) {
	// Check if the channel exists
	if !m.ChannelExists(channel) {
		m.logger.Info(
			"info: channel does not exist",
			zap.String("channel", channel),
		)

		// Create a new channel
		if err := m.CreateChannel(channel); err != nil {
			return 0, ErrInternal
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Append the message to the channel
	m.data[channel] = append(m.data[channel], message)

	// Return the index of the message in the channel
	return int64(len(m.data[channel]) - 1), nil
}

// GetMessages retrieves all messages from the specified channel
func (m *MemoryStorage) GetMessages(channel string, offset int64) ([]interface{}, int64, error) {
	if offset < -1 {
		return nil, 0, ErrInvalidOffset
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, exists := m.data[channel]
	if !exists {
		m.logger.Info(
			"info: channel does not exist",
			zap.String("channel", channel),
		)
		return nil, 0, fmt.Errorf("channel '%s' does not exist", channel)
	}

	// Check if the offset is greater than the number of messages
	if int(offset) >= len(messages) {
		return nil, 0, ErrInvalidOffset
	}

	// Return only the latest message if the offset is -1
	if offset == -1 {
		return []interface{}(nil), int64(len(messages) - 1), nil
	}

	// Limit the number of messages to be returned
	var endOffset int
	if int(offset)+m.batchSize > len(messages) {
		endOffset = len(messages)
	} else {
		endOffset = int(offset) + m.batchSize
	}

	// Return the copy of the messages to prevent data mutation
	return append([]interface{}(nil), messages[offset:endOffset]...), int64(endOffset - 1), nil
}

// CreateChannel creates a new channel
func (m *MemoryStorage) CreateChannel(channel string) error {
	if m.ChannelExists(channel) {
		m.logger.Info(
			"info: channel already exists",
			zap.String("channel", channel),
		)
		return fmt.Errorf("channel '%s' already exists", channel)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[channel] = []interface{}{}
	return nil
}

// ChannelExists checks if a channel exists
func (m *MemoryStorage) ChannelExists(channel string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[channel]
	return exists
}

// DeleteChannel deletes an existing channel
func (m *MemoryStorage) DeleteChannel(channel string) error {
	if !m.ChannelExists(channel) {
		m.logger.Info(
			"info: channel does not exist",
			zap.String("channel", channel),
		)
		return fmt.Errorf("channel '%s' does not exist", channel)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, channel)
	return nil
}
