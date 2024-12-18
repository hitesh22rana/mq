// pkg/storage/memory_storage.go

package storage

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// chunk represents a chunk of data
type chunk struct {
	data interface{}
	prev *chunk
	next *chunk
}

// chunkList represents a linked list of chunks
type chunkList struct {
	head *chunk
	tail *chunk
	len  uint64
}

// MemoryStorageOptions represents the options for the MemoryStorage
type MemoryStorageOptions struct {
	BatchSize uint64
}

// MemoryStorage is an in-memory implementation of the Storage interface
type MemoryStorage struct {
	mu                       sync.RWMutex
	logger                   *zap.Logger
	batchSize                uint64
	data                     map[string]*chunkList
	subscriberToChannelChunk map[string]map[string]*chunk
}

// NewMemoryStorage initializes a new MemoryStorage instance
func NewMemoryStorage(logger *zap.Logger, options *MemoryStorageOptions) *MemoryStorage {
	return &MemoryStorage{
		mu:                       sync.RWMutex{},
		logger:                   logger,
		batchSize:                options.BatchSize,
		data:                     make(map[string]*chunkList),
		subscriberToChannelChunk: make(map[string]map[string]*chunk),
	}
}

// SaveMessage saves a message to the specified channel
func (m *MemoryStorage) SaveMessage(channel string, message interface{}) (uint64, error) {
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

	// Get the list of messages in the channel
	msgList := m.data[channel]

	// Make a new chunk
	chunk := &chunk{
		data: message,
		prev: nil,
		next: nil,
	}

	// Append the message to the channel
	if msgList.head == nil {
		msgList.head = chunk
		msgList.tail = chunk
		msgList.len = 1
	} else {
		chunk.prev = msgList.tail
		msgList.tail.next = chunk
		msgList.tail = msgList.tail.next
		msgList.len++
	}

	// Return the index of the message in the channel
	return msgList.len, nil
}

// GetMessages retrieves all messages from the specified channel
func (m *MemoryStorage) GetMessages(channel string, subscriberID string, offset uint64, isFromLatest *bool) ([]interface{}, uint64, error) {
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
	if offset >= messages.len {
		return []interface{}(nil), 0, ErrInvalidOffset
	}

	// Initialize the subscriberToChannelChunk map if it does not exist
	if _, exists := m.subscriberToChannelChunk[subscriberID]; !exists {
		m.subscriberToChannelChunk[subscriberID] = make(map[string]*chunk)
	}

	// Return only the latest message if isStart is false
	if *isFromLatest {
		*isFromLatest = false

		// Update the last chunk in the subscriberToChannelChunk map to the latest message
		m.subscriberToChannelChunk[subscriberID][channel] = messages.tail
		return []interface{}(nil), messages.len - 1, nil
	}

	// Limit the number of messages to be returned
	var endIndx uint64 = 0
	if offset+m.batchSize > messages.len {
		endIndx = messages.len
	} else {
		endIndx = offset + m.batchSize
	}

	// Copy the messages from the channel
	data := make([]interface{}, 0)

	// Get the iterator to the start offset
	iterator := messages.head
	prevChunk := m.subscriberToChannelChunk[subscriberID][channel]
	if prevChunk != nil {
		iterator = prevChunk.next
	}

	for i := offset; i < endIndx; i++ {
		data = append(data, iterator.data)
		if iterator.next == nil {
			break
		}

		iterator = iterator.next
	}

	// Update the last chunk in the subscriberToChannelChunk map
	m.subscriberToChannelChunk[subscriberID][channel] = iterator

	// Return the messages and the next offset
	return data, endIndx - 1, nil
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

	m.data[channel] = &chunkList{
		head: nil,
		tail: nil,
		len:  0,
	}

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

	// Remove the channel from the subscriberToChannelChunk map
	for ch := range m.subscriberToChannelChunk[channel] {
		delete(m.subscriberToChannelChunk[channel], ch)
	}

	return nil
}
