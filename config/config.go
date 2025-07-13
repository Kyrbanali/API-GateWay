package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Cache    CacheConfig
}

type AppConfig struct {
	Port string
}

type PostgresConfig struct {
	Host string
	User string
	Pass string
	DB   string
	Port string
}

type CacheConfig struct {
	TTL             string
	CleanupInterval string
}

func Load() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "read config.yaml")
	}

	viper.AutomaticEnv()
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return &cfg, nil
}

func (p *PostgresConfig) BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		p.User,
		p.Pass,
		p.Host,
		p.Port,
		p.DB,
	)
}
