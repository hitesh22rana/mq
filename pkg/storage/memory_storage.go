// pkg/storage/memory_storage.go

package storage

import (
	"fmt"
	"io"
	"sync"

	"github.com/rosedblabs/wal"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// chunk represents a chunk of data
type chunk struct {
	data *pb.Message
	prev *chunk
	next *chunk
}

// chunkList represents a linked list of chunks
type chunkList struct {
	head *chunk
	tail *chunk
	len  uint64
}

// appendChunk appends a chunk to the chunk list
func (cl *chunkList) appendChunk(chunk *chunk) {
	// Append the message to the list
	if cl.head == nil {
		cl.head = chunk
		cl.tail = chunk
		cl.len = 1
	} else {
		chunk.prev = cl.tail
		cl.tail.next = chunk
		cl.tail = cl.tail.next
		cl.len++
	}
}

// MemoryStorageOptions represents the options for the MemoryStorage
type MemoryStorageOptions struct {
	Wal           *wal.WAL
	BatchSize     uint64
	SyncOnStartup bool
}

// MemoryStorage is an in-memory implementation of the Storage interface
type MemoryStorage struct {
	mu                       sync.RWMutex
	logger                   *zap.Logger
	wal                      *wal.WAL
	batchSize                uint64
	data                     map[string]*chunkList
	subscriberToChannelChunk map[string]map[string]*chunk
}

// NewMemoryStorage initializes a new MemoryStorage instance
func NewMemoryStorage(logger *zap.Logger, options *MemoryStorageOptions) *MemoryStorage {
	m := &MemoryStorage{
		mu:                       sync.RWMutex{},
		logger:                   logger,
		wal:                      options.Wal,
		batchSize:                options.BatchSize,
		data:                     make(map[string]*chunkList),
		subscriberToChannelChunk: make(map[string]map[string]*chunk),
	}

	if !options.SyncOnStartup {
		return m
	}

	// Inform the user that the storage is being synced
	m.logger.Info("info: syncing storage on startup, this may take a while")

	// Load data from the Write-Ahead Log (WAL)
	reader := m.wal.NewReader()
	for {
		data, _, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			m.logger.Fatal(
				"fatal: failed to read WAL",
				zap.Error(err),
			)
		}

		// Unmarshal the protobuf data
		entry := &pb.WalEntry{}
		if err := proto.Unmarshal(data, entry); err != nil {
			m.logger.Error(
				"error: failed to unmarshal data",
				zap.Error(err),
			)
			break
		}

		channel := entry.GetChannel()
		message := entry.GetMessage()

		// Create the channel if it does not exist
		if !m.ChannelExists(channel) {
			_ = m.CreateChannel(channel)
			m.logger.Info(
				"info: created channel",
				zap.String("channel", channel),
			)
		}

		// Get the list of messages in the channel
		msgList := m.data[channel]

		// Make a new chunk and append it to the list
		msgList.appendChunk(
			&chunk{
				data: message,
				prev: nil,
				next: nil,
			},
		)
	}

	// Inform the user that the storage has been synced
	m.logger.Info("info: storage synced successfully")
	return m
}

// SaveMessage saves a message to the specified channel
func (m *MemoryStorage) SaveMessage(
	channel string,
	message *pb.Message,
) (uint64, error) {
	// Create the channel if it does not exist
	if !m.ChannelExists(channel) {
		_ = m.CreateChannel(channel)
		m.logger.Info(
			"info: created channel",
			zap.String("channel", channel),
		)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Get the list of messages in the channel
	msgList := m.data[channel]

	// Write the message to the Write-Ahead Log (WAL)
	entry := &pb.WalEntry{
		Channel: channel,
		Message: message,
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		m.logger.Error(
			"error: failed to marshal data",
			zap.Error(err),
		)

		return msgList.len, ErrInternal
	}

	// Write the data to the WAL
	if _, err = m.wal.Write(data); err != nil {
		m.logger.Error(
			"error: failed to write to WAL",
			zap.Error(err),
		)

		return msgList.len, ErrInternal
	}

	// Make a new chunk and append it to the list
	msgList.appendChunk(
		&chunk{
			data: message,
			prev: nil,
			next: nil,
		},
	)

	// Return the index of the message in the channel
	return msgList.len, nil
}

// GetMessages retrieves all messages from the specified channel
func (m *MemoryStorage) GetMessages(
	channel string,
	subscriberID string,
	offset uint64,
) ([]*pb.Message, uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, exists := m.data[channel]
	if !exists {
		m.logger.Error(
			"error: channel does not exist",
			zap.String("channel", channel),
		)
		return nil, 0, fmt.Errorf("channel '%s' does not exist", channel)
	}

	// Initialize the subscriberToChannelChunk map if it does not exist
	if _, exists := m.subscriberToChannelChunk[subscriberID]; !exists {
		m.subscriberToChannelChunk[subscriberID] = make(map[string]*chunk)
	}

	// Check if the offset is valid
	if offset >= messages.len {
		// Return only the latest message if the offset is set to the latest
		if offset == OffsetLatest {
			// Update the last chunk in the subscriberToChannelChunk map to the latest message
			m.subscriberToChannelChunk[subscriberID][channel] = messages.tail
			return []*pb.Message(nil), messages.len - 1, nil
		}

		return []*pb.Message(nil), 0, ErrInvalidOffset
	}

	// Limit the number of messages to be returned
	var endIndx uint64 = 0
	if offset+m.batchSize > messages.len {
		endIndx = messages.len
	} else {
		endIndx = offset + m.batchSize
	}

	// Copy the messages from the channel
	data := make([]*pb.Message, 0)

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

// RemoveChannelFromSubscriberMap removes the channel from the subscriberToChannelChunk map
func (m *MemoryStorage) RemoveChannelFromSubscriberMap(
	channel string,
	subscriberID string,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove the channel from the subscriberToChannelChunk map
	delete(m.subscriberToChannelChunk[subscriberID], channel)
}
