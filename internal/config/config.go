package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	DB   DBConfig
}

type AppConfig struct {
	JWTSecret []byte `envconfig:"JWT_SECRET"`
}

type HTTPConfig struct {
	Host string `envconfig:"HTTP_HOST"`
	Port int    `envconfig:"HTTP_PORT"`
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	Name     string `envconfig:"DB_NAME"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process env vars: %w", err)
	}

	return &cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s pool_max_conns=1500 pool_max_conn_lifetime=1m",
		c.Host, c.Port, c.Name, c.User, c.Password,
	)
}

func (c *HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
