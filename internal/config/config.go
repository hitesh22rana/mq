package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = ""

type Configuration struct {
	Storage
	Broker
	Producer
	Subscriber
	Environment
}

type Storage struct {
	MemoryStorageBatchSize int `envconfig:"MEMORY_STORAGE_BATCH_SIZE" default:"500"`
}

type Broker struct {
	BrokerPort                    string        `envconfig:"BROKER_PORT" required:"true"`
	BrokerURL                     string        `envconfig:"BROKER_URL" required:"true"`
	BrokerGracefulShutdownTimeout time.Duration `envconfig:"BROKER_GRACEFUL_SHUTDOWN_TIMEOUT" default:"5s"`
}

type Producer struct {
	KeepAliveTime       time.Duration `envconfig:"KEEPALIVE_TIME" default:"10s"`
	KeepAliveTimeout    time.Duration `envconfig:"KEEPALIVE_TIMEOUT" default:"5s"`
	PermitWithoutStream bool          `envconfig:"PERMIT_WITHOUT_STREAM" default:"true"`
}

type Subscriber struct {
	// DataPullingInterval is the interval at which the subscriber pulls data from the broker (in milliseconds)
	DataPullingInterval int64 `envconfig:"DATA_PULLING_INTERVAL" default:"-1"`
}

type Environment struct {
	Env string `envconfig:"ENV" default:"development"`
}

func Load() (*Configuration, error) {
	_ = godotenv.Overload()

	var cfg Configuration
	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
