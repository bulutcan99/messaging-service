package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "config.yml"
)

type (
	Container struct {
		App      `yaml:"app"`
		Log      `yaml:"log"`
		HTTP     *HTTP     `yaml:"http"`
		Token    *Token    `yaml:"token"`
		MSSQL    *MSSQL    `yaml:"mssql"`
		Settings *Settings `yaml:"settings"`
	}

	App struct {
		Name string `env-required:"true" yaml:"name" env:"APP_NAME"`
	}

	Log struct {
		Level int `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Token struct {
		Secret string `env-required:"true" yaml:"secret" env:"TOKEN_SECRET"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port int    `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	MSSQL struct {
		URL string `env-required:"true" yaml:"url" env:"MSSQL_URL"`
	}

	Settings struct {
		ServerReadTimeout int `env-required:"true" yaml:"server_read_timeout" env:"SERVER_READ_TIMEOUT"`
		MSSQLConnAttempts int `env-required:"true" yaml:"mssql_conn_attempts" env:"MSSQL_CONN_ATTEMPTS"`
		MSSQLConnTimeout  int `env-required:"true" yaml:"mssql_conn_timeout" env:"MSSQL_CONN_TIMEOUT"`
	}
)

func NewConfig() (*Container, error) {
	cfg := new(Container)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
