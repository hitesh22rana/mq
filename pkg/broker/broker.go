// pkg/broker/broker.go

package broker

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"github.com/hitesh22rana/mq/pkg/storage"
	"github.com/hitesh22rana/mq/pkg/utils"
)

var (
	// ErrFailedToSaveMessage is returned when the broker fails to save a message
	ErrFailedToSaveMessage = errors.New("error: failed to save message")

	// ErrChannelDoesNotExist is returned when the broker tries to publish a message to a non-existent channel
	ErrChannelDoesNotExist = errors.New("error: channel does not exist")

	// ErrChannelAlreadyExists is returned when the broker tries to create a channel that already exists
	ErrChannelAlreadyExists = errors.New("error: channel already exists")

	// ErrSubscriberDoesNotExist is returned when the broker tries to unsubscribe a subscriber from a channel in which the subscriber does not exist
	ErrSubscriberDoesNotExist = errors.New("error: subscriber is not subscribed to the channel")
)

// Broker defines the interface for the message broker
type Broker interface {
	createChannel(context.Context, string) error
	publish(context.Context, string, *message) error
	subscribe(context.Context, string, string, chan<- *message) error
	unsubscribe(context.Context, string, string) error
}

type message struct {
	id        string
	payload   string
	timeStamp int64
}

// Service is the implementation of the Broker interface
type Service struct {
	mu          sync.RWMutex
	logger      *zap.Logger
	storage     storage.Storage
	subscribers map[string]map[string]chan<- *message
}

// NewService returns a new broker service
func NewService(logger *zap.Logger, storage storage.Storage) *Service {
	return &Service{
		mu:          sync.RWMutex{},
		logger:      logger,
		storage:     storage,
		subscribers: make(map[string]map[string]chan<- *message),
	}
}

// Server is the broker service implementation for gRPC
type Server struct {
	logger    *zap.Logger
	validator utils.Validator
	generator utils.Generator
	srv       Broker
}

// NewServer returns a new broker server
func NewServer(logger *zap.Logger, srv *Service) *Server {
	return &Server{
		logger:    logger,
		validator: utils.NewValidator(),
		generator: utils.NewGenerator(),
		srv:       srv,
	}
}
