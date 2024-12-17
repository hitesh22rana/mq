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
	Broker
	Publisher
	Subscriber
	Environment
}

// Storage represents the configuration for the storage mechanism
type Storage struct {
	MemoryStorageBatchSize int `envconfig:"MEMORY_STORAGE_BATCH_SIZE" default:"500"`
}

// Broker represents the configuration for the broker
type Broker struct {
	BrokerPort                    int           `envconfig:"BROKER_PORT" required:"true"`
	BrokerHost                    string        `envconfig:"BROKER_HOST" required:"true"`
	BrokerGracefulShutdownTimeout time.Duration `envconfig:"BROKER_GRACEFUL_SHUTDOWN_TIMEOUT" default:"5s"`
}

// Publisher represents the configuration for the publisher
type Publisher struct {
	KeepAliveTime       time.Duration `envconfig:"KEEPALIVE_TIME" default:"10s"`
	KeepAliveTimeout    time.Duration `envconfig:"KEEPALIVE_TIMEOUT" default:"5s"`
	PermitWithoutStream bool          `envconfig:"PERMIT_WITHOUT_STREAM" default:"true"`
}

// Subscriber represents the configuration for the subscriber
type Subscriber struct {
	// DataPullingInterval is the interval at which the subscriber pulls data from the broker (in milliseconds)
	DataPullingInterval uint64 `envconfig:"DATA_PULLING_INTERVAL" default:"100"`
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
