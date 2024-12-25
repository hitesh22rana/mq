// pkg/mq/mq.go

package mq

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
	"github.com/hitesh22rana/mq/pkg/storage"
	"github.com/hitesh22rana/mq/pkg/utils"
)

var (
	// ErrFailedToSaveMessage is returned when the mq fails to save a message
	ErrFailedToSaveMessage = errors.New("error: failed to save message")

	// ErrUnableToCreateChannel is returned when the mq fails to create a channel
	ErrUnableToCreateChannel = errors.New("error: unable to create channel")

	// ErrChannelDoesNotExist is returned when the mq tries to publish a message to a non-existent channel
	ErrChannelDoesNotExist = errors.New("error: channel does not exist")

	// ErrChannelAlreadyExists is returned when the mq tries to create a channel that already exists
	ErrChannelAlreadyExists = errors.New("error: channel already exists")

	// ErrSubscriberDoesNotExist is returned when the mq tries to unsubscribe a subscriber from a channel in which the subscriber does not exist
	ErrSubscriberDoesNotExist = errors.New("error: subscriber is not subscribed to the channel")
)

// MQ defines the interface for the mq
type MQ interface {
	createChannel(context.Context, channel) error
	publish(context.Context, channel, *pb.Message) error
	subscribe(context.Context, *pb.Subscriber, pb.Offset, uint64, channel, chan<- *pb.Message) error
	unsubscribe(context.Context, *pb.Subscriber, channel) error
}

// Channel represents a message channel
type channel string

// Service is the implementation of the MQ interface
type Service struct {
	mu                   sync.RWMutex
	logger               *zap.Logger
	storage              storage.Storage
	channelToSubscribers map[channel]map[*pb.Subscriber]struct{}
}

// ServiceOptions represents the options for the mq service
type ServiceOptions struct {
	Storage storage.Storage
}

// NewService returns a new mq service
func NewService(logger *zap.Logger, options *ServiceOptions) *Service {
	return &Service{
		mu:                   sync.RWMutex{},
		logger:               logger,
		storage:              options.Storage,
		channelToSubscribers: make(map[channel]map[*pb.Subscriber]struct{}),
	}
}

// Server is the mq service implementation for gRPC
type Server struct {
	logger    *zap.Logger
	validator utils.Validator
	generator utils.Generator
	srv       MQ
}

// ServerOptions represents the options for the mq server
type ServerOptions struct {
	Validator utils.Validator
	Generator utils.Generator
	Service   MQ
}

// NewServer returns a new mq server
func NewServer(logger *zap.Logger, options *ServerOptions) *Server {
	return &Server{
		logger:    logger,
		validator: options.Validator,
		generator: options.Generator,
		srv:       options.Service,
	}
}
