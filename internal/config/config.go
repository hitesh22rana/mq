package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = ""

type Configuration struct {
	Broker
	Producer
	Subscriber
	Environment
}

type Broker struct {
	BrokerPort string `envconfig:"BROKER_PORT" required:"true"`
	BrokerURL  string `envconfig:"BROKER_URL" required:"true"`
}

type Producer struct {
	KeepAliveTime       time.Duration `envconfig:"KEEPALIVE_TIME" default:"10s"`
	KeepAliveTimeout    time.Duration `envconfig:"KEEPALIVE_TIMEOUT" default:"5s"`
	PermitWithoutStream bool          `envconfig:"PERMIT_WITHOUT_STREAM" default:"true"`
}

type Subscriber struct {
	KeepAliveTime       time.Duration `envconfig:"KEEPALIVE_TIME" default:"10s"`
	KeepAliveTimeout    time.Duration `envconfig:"KEEPALIVE_TIMEOUT" default:"5s"`
	PermitWithoutStream bool          `envconfig:"PERMIT_WITHOUT_STREAM" default:"true"`
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
