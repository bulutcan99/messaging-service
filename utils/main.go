package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel        = "info"
	DefaultHost            = "localhost"
	DefaultPort            = 8080
	DefaultDatabaseURI     = "mongodb://localhost:27017"
	DefaultDatabaseName    = "messaging_service"
	DefaultUsersCollection = "users"
	DefaultCacheURI        = "redis://localhost:6379"
)

type LogLevel string

const (
	DebugLevel   LogLevel = "debug"
	InfoLevel    LogLevel = "info"
	WarningLevel LogLevel = "warning"
	ErrorLevel   LogLevel = "error"
	FatalLevel   LogLevel = "fatal"
)

type Config struct {
	LogLevel        LogLevel `mapstructure:"LogLevel" validate:"oneof=debug info warn error fatal"`
	DatabaseURI     string   `mapstructure:"DatabaseURI"`
	DatabaseName    string   `mapstructure:"DatabaseName"`
	UsersCollection string   `mapstructure:"UsersCollection"`
	CacheURI        string   `mapstructure:"CacheURI"`
	Host            string   `mapstructure:"Host"`
	Port            int      `mapstructure:"Port"`
	JwtSecretKey    string   `mapstructure:"JwtSecretKey" validate:"required"`
	JwtPublicKey    string   `mapstructure:"JwtPublicKey" validate:"required"`
}

func (c *Config) setDefaults() {
	if c.LogLevel == "" {
		c.LogLevel = DefaultLogLevel
		fmt.Printf("LogLevel not set, using default: %s", DefaultLogLevel)
	}
	if c.Host == "" {
		c.Host = DefaultHost
		fmt.Printf("Host not set, using default: %s", DefaultHost)
	}
	if c.Port == 0 {
		c.Port = DefaultPort
		fmt.Printf("Port not set, using default: %d", DefaultPort)
	}
	if c.DatabaseURI == "" {
		c.DatabaseURI = DefaultDatabaseURI
		fmt.Printf("DatabaseURI not set, using default: %s", DefaultDatabaseURI)
	}
	if c.DatabaseName == "" {
		c.DatabaseName = DefaultDatabaseName
		fmt.Printf("DatabaseName not set, using default: %s", DefaultDatabaseName)
	}
	if c.UsersCollection == "" {
		c.UsersCollection = DefaultUsersCollection
		fmt.Printf("UsersCollection not set, using default: %s", DefaultUsersCollection)
	}
	if c.CacheURI == "" {
		c.CacheURI = DefaultCacheURI
		fmt.Printf("CacheURI not set, using default: %s", DefaultCacheURI)
	}
}

func (c *Config) validate() error {
	if c.JwtSecretKey == "" {
		return fmt.Errorf("JwtSecretKey is required")
	}
	if c.JwtPublicKey == "" {
		return fmt.Errorf("JwtPublicKey is required")
	}

	validLogLevels := map[LogLevel]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	return nil
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults and log warnings for empty values
	config.setDefaults()

	// Validate required fields
	if err := config.validate(); err != nil {
		return config, err
	}

	return config, nil
}
