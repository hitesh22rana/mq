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

// Storage represents the configuration for the storage mechanism
type Storage struct {
	StorageBatchSize uint64 `envconfig:"STORAGE_BATCH_SIZE" default:"500"`
}

// Wal represents the configuration for the Write-Ahead Log (WAL)
type Wal struct {
	WalDirPath        string `envconfig:"WAL_DIR_PATH" default:"/tmp"`
	WalSegmentSize    int64  `envconfig:"WAL_SEGMENT_SIZE" default:"52428800"` // 50MB (5,24,28,800) bytes
	WalSegmentFileExt string `envconfig:"WAL_SEGMENT_FILE_EXT" default:".wal"`
	WalSync           bool   `envconfig:"WAL_SYNC" default:"false"`
	WalBytesPerSync   uint32 `envconfig:"WAL_BYTES_PER_SYNC" default:"0"`
}

// Broker represents the configuration for the broker
type Broker struct {
	BrokerPort                    int           `envconfig:"BROKER_PORT" required:"true"`
	BrokerHost                    string        `envconfig:"BROKER_HOST" required:"true"`
	BrokerGracefulShutdownTimeout time.Duration `envconfig:"BROKER_GRACEFUL_SHUTDOWN_TIMEOUT" default:"5s"`
}

// GrpcServer represents the configuration for the gRPC server
type GrpcServer struct {
	GrpcServerMaxRecvMsgSize int `envconfig:"GRPC_SERVER_MAX_RECV_MSG_SIZE" default:"4194304"` // 4 MB (41,94,304) bytes
}

// Publisher represents the configuration for the publisher
type Publisher struct {
	PublisherKeepAliveTime       time.Duration `envconfig:"PUBLISHER_KEEPALIVE_TIME" default:"10s"`
	PublisherKeepAliveTimeout    time.Duration `envconfig:"PUBLISHER_KEEPALIVE_TIMEOUT" default:"5s"`
	PublisherPermitWithoutStream bool          `envconfig:"PUBLISHER_PERMIT_WITHOUT_STREAM" default:"true"`
}

// Subscriber represents the configuration for the subscriber
type Subscriber struct {
	// SubscriberDataPullingInterval is the interval at which the subscriber pulls data from the broker (in milliseconds)
	SubscriberDataPullingInterval uint64 `envconfig:"SUBSCRIBER_DATA_PULLING_INTERVAL" default:"100"`
}

// Environment represents the environment configuration
type Environment struct {
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
