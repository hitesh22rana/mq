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
	ErrFailedToSaveMessage = errors.New("failed to save message")

	// ErrChannelDoesNotExist is returned when the broker tries to publish a message to a non-existent channel
	ErrChannelDoesNotExist = errors.New("channel does not exist")

	// ErrSubscriberAlreadyExists is returned when the broker tries to subscribe a subscriber that already exists
	ErrSubscriberAlreadyExists = errors.New("subscriber already exists")

	// ErrSubscriberDoesNotExist is returned when the broker tries to unsubscribe a subscriber from a channel in which the subscriber does not exist
	ErrSubscriberDoesNotExist = errors.New("subscriber is not subscribed to the channel")
)

// Broker defines the interface for the message broker
type Broker interface {
	publish(context.Context, string, message) error
	subscribe(context.Context, string, string, chan<- *message) error
	unsubscribe(context.Context)
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
func NewService(storage storage.Storage, logger *zap.Logger) *Service {
	return &Service{
		mu:          sync.RWMutex{},
		logger:      logger,
		storage:     storage,
		subscribers: make(map[string]map[string]chan<- *message),
	}
}

// Server is the broker service implementation for gRPC
type Server struct {
	validator utils.Validator
	generator utils.Generator
	srv       *Service
}

// NewServer returns a new broker server
func NewServer(srv *Service) *Server {
	return &Server{
		validator: utils.NewValidator(),
		generator: utils.NewGenerator(),
		srv:       srv,
	}
}
