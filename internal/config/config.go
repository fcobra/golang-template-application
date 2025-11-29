package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP        HTTPConfig        `yaml:"http"`
	Auth        AuthConfig        `yaml:"auth"`
	Logger      LoggerConfig      `yaml:"logger"`
	Postgres    PostgresConfig    `yaml:"postgres"`
	Redis       RedisConfig       `yaml:"redis"`
	Pushgateway PushgatewayConfig `yaml:"pushgateway"`
	Sentry      SentryConfig      `yaml:"sentry"`
}

type AuthConfig struct {
	Provider string `yaml:"provider" env-default:"inmemory"`
}

type HTTPConfig struct {
	Host         string        `yaml:"host" env:"HTTP_HOST"`
	Port         string        `yaml:"port" env:"HTTP_PORT"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type LoggerConfig struct {
	Enabled     bool   `yaml:"enabled" env-default:"true"`
	Level       string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Destination string `yaml:"destination" env-default:"stdout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"PG_HOST"`
	Port     string `yaml:"port" env:"PG_PORT"`
	User     string `yaml:"user" env:"PG_USER"`
	Password string `yaml:"password" env:"PG_PASSWORD"`
	DBName   string `yaml:"dbname" env:"PG_DBNAME"`
	SSLMode  string `yaml:"sslmode" env:"PG_SSLMODE"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

type RedisConfig struct {
	Enabled  bool   `yaml:"enabled" env-default:"true"`
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Port     string `yaml:"port" env:"REDIS_PORT"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

type PushgatewayConfig struct {
	Enabled bool   `yaml:"enabled" env-default:"true"`
	URL     string `yaml:"url" env:"PUSHGATEWAY_URL"`
}

type SentryConfig struct {
	Enabled bool   `yaml:"enabled" env-default:"false"`
	Dsn     string `yaml:"dsn" env:"SENTRY_DSN"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file does not exist: %s", configPath))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("cannot read config: %s", err))
	}

	localConfigPath := strings.TrimSuffix(configPath, filepath.Ext(configPath)) + ".local.yaml"
	if _, err := os.Stat(localConfigPath); err == nil {
		if err := cleanenv.ReadConfig(localConfigPath, &cfg); err != nil {
			panic(fmt.Sprintf("cannot read local config: %s", err))
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("cannot read environment variables: %s", err))
	}

	return &cfg
}