// pkg/broker/broker.go

package broker

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"github.com/hitesh22rana/mq/pkg/proto/event"
	"github.com/hitesh22rana/mq/pkg/storage"
	"github.com/hitesh22rana/mq/pkg/utils"
)

var (
	// ErrFailedToSaveMessage is returned when the broker fails to save a message
	ErrFailedToSaveMessage = errors.New("error: failed to save message")

	// ErrUnableToCreateChannel is returned when the broker fails to create a channel
	ErrUnableToCreateChannel = errors.New("error: unable to create channel")

	// ErrChannelDoesNotExist is returned when the broker tries to publish a message to a non-existent channel
	ErrChannelDoesNotExist = errors.New("error: channel does not exist")

	// ErrChannelAlreadyExists is returned when the broker tries to create a channel that already exists
	ErrChannelAlreadyExists = errors.New("error: channel already exists")

	// ErrSubscriberDoesNotExist is returned when the broker tries to unsubscribe a subscriber from a channel in which the subscriber does not exist
	ErrSubscriberDoesNotExist = errors.New("error: subscriber is not subscribed to the channel")
)

// Broker defines the interface for the message broker
type Broker interface {
	createChannel(context.Context, channel) error
	publish(context.Context, channel, *event.Message) error
	subscribe(context.Context, *event.Subscriber, event.Offset, uint64, channel, chan<- *event.Message) error
	unsubscribe(context.Context, *event.Subscriber, channel) error
}

// Channel represents a message channel
type channel string

// Service is the implementation of the Broker interface
type Service struct {
	mu                   sync.RWMutex
	logger               *zap.Logger
	storage              storage.Storage
	channelToSubscribers map[channel]map[*event.Subscriber]struct{}
}

// ServiceOptions represents the options for the broker service
type ServiceOptions struct {
	Storage storage.Storage
}

// NewService returns a new broker service
func NewService(logger *zap.Logger, options *ServiceOptions) *Service {
	return &Service{
		mu:                   sync.RWMutex{},
		logger:               logger,
		storage:              options.Storage,
		channelToSubscribers: make(map[channel]map[*event.Subscriber]struct{}),
	}
}

// Server is the broker service implementation for gRPC
type Server struct {
	logger    *zap.Logger
	validator utils.Validator
	generator utils.Generator
	srv       Broker
}

// ServerOptions represents the options for the broker server
type ServerOptions struct {
	Validator utils.Validator
	Generator utils.Generator
	Service   Broker
}

// NewServer returns a new broker server
func NewServer(logger *zap.Logger, options *ServerOptions) *Server {
	return &Server{
		logger:    logger,
		validator: options.Validator,
		generator: options.Generator,
		srv:       options.Service,
	}
}
