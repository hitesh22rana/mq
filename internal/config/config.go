// internal/config/config.go

package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// envPrefix is the prefix used for environment variables
const envPrefix = ""

// Configuration represents the configuration of the application
type Configuration struct {
	Storage
	Wal
	Broker
	GrpcServer
	Publisher
	Subscriber
	Environment
}

// Storage holds configuration settings for storage operations.
type Storage struct {
	// StorageBatchSize specifies the batch size for storage operations.
	// default: 500
	StorageBatchSize uint64 `envconfig:"STORAGE_BATCH_SIZE" default:"500"`

	// StorageSyncOnStartup indicates whether to perform a storage synchronization on startup.
	// default: true
	StorageSyncOnStartup bool `envconfig:"STORAGE_SYNC_ON_STARTUP" default:"true"`
}

// Wal represents the configuration for the Write-Ahead Logging (WAL) mechanism.
type Wal struct {
	// WalDirPath specifies the directory path where the WAL segments are stored
	// default: ./data
	WalDirPath string `envconfig:"WAL_DIR_PATH" default:"./data"`

	// WalSegmentSize defines the size (in bytes) for each WAL segment
	// default: 52428800 bytes (50MB)
	WalSegmentSize int64 `envconfig:"WAL_SEGMENT_SIZE" default:"52428800"` // 50MB (5,24,28,800) bytes

	// WalSegmentFileExt determines the file extension used for WAL segments.
	// default: .wal
	WalSegmentFileExt string `envconfig:"WAL_SEGMENT_FILE_EXT" default:".wal"`

	// WalSync controls whether each write operation should be fsync’d to ensure strong consistency.
	// default: true
	WalSync bool `envconfig:"WAL_SYNC" default:"true"`

	// WalBytesPerSync specifies the number of bytes to write before calling fsync.
	// default: 5000 bytes (5KB)
	WalBytesPerSync uint32 `envconfig:"WAL_BYTES_PER_SYNC" default:"5000"`
}

// Broker represents the configuration settings for a message broker.
type Broker struct {
	// BrokerPort specifies the port on which the broker listens. This field is required.
	BrokerPort int `envconfig:"BROKER_PORT" required:"true"`

	// BrokerHost specifies the host address of the broker. This field is required.
	BrokerHost string `envconfig:"BROKER_HOST" required:"true"`

	// BrokerGracefulShutdownTimeout specifies the duration to wait for a graceful shutdown
	// default: 5s
	BrokerGracefulShutdownTimeout time.Duration `envconfig:"BROKER_GRACEFUL_SHUTDOWN_TIMEOUT" default:"5s"`
}

// GrpcServer holds the configuration settings for a gRPC server.
type GrpcServer struct {
	// GrpcServerMaxRecvMsgSize specifies the maximum size of a message that the server can receive.
	// default: 4194304 (4 MB)
	GrpcServerMaxRecvMsgSize int `envconfig:"GRPC_SERVER_MAX_RECV_MSG_SIZE" default:"4194304"` // 4 MB (41,94,304) bytes
}

// Publisher holds the configuration settings for a publisher.
type Publisher struct {
	// PublisherKeepAliveTime is the duration for the keep-alive time (interval between keep-alive pings).
	// default: 10s
	PublisherKeepAliveTime time.Duration `envconfig:"PUBLISHER_KEEPALIVE_TIME" default:"10s"`

	// PublisherKeepAliveTimeout is the duration for the keep-alive timeout (time to wait for a keep-alive response before considering the connection dead)
	// default: 5s
	PublisherKeepAliveTimeout time.Duration `envconfig:"PUBLISHER_KEEPALIVE_TIMEOUT" default:"5s"`

	// PublisherPermitWithoutStream is a boolean indicating whether to allow keep-alive pings when there are no active streams.
	// default: true
	PublisherPermitWithoutStream bool `envconfig:"PUBLISHER_PERMIT_WITHOUT_STREAM" default:"true"`
}

// Subscriber represents the configuration for the subscriber
type Subscriber struct {
	// SubscriberDataPullingInterval is the interval at which the subscriber pulls data from the broker
	// default: 100 (100ms)
	SubscriberDataPullingInterval uint64 `envconfig:"SUBSCRIBER_DATA_PULLING_INTERVAL" default:"100"`
}

// Environment holds the configuration for the application's environment settings.
type Environment struct {
	// Env represents the current environment
	// default:  development
	Env string `envconfig:"ENV" default:"development"`
}

// Load loads the configuration from the environment
func Load() (*Configuration, error) {
	_ = godotenv.Overload()

	var cfg Configuration
	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
