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
	Server
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

// Server holds the configuration settings for server.
type Server struct {
	// ServerPort specifies the port on which the server listens. This field is required.
	ServerPort int `envconfig:"SERVER_PORT" required:"true"`

	// ServerGracefulShutdownTimeout specifies the duration to wait for a graceful shutdown
	// default: 5s
	ServerGracefulShutdownTimeout time.Duration `envconfig:"SERVER_GRACEFUL_SHUTDOWN_TIMEOUT" default:"5s"`

	// ServerMaxRecvMsgSize specifies the maximum size of a message that the server can receive.
	// default: 4194304 (4 MB)
	ServerMaxRecvMsgSize int `envconfig:"SERVER_MAX_RECV_MSG_SIZE" default:"4194304"` // 4 MB (41,94,304) bytes
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
