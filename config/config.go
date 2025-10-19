package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/Kyrbanali/API-GateWay/internal/worker"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port string
	}
	Postgres struct {
		Host string
		User string
		Pass string
		DB   string
		Port string
	}
	Cache struct {
		TTL             time.Duration
		CleanupInterval time.Duration
	}
	Workers worker.Config

	Grafana struct {
		User string
		Pass string
	}
}

//go:embed config.yaml
var defaultYAMLConfig []byte

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(defaultYAMLConfig)); err != nil {
		return nil, errors.Wrap(err, "read embedded config")
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return &cfg, nil
}

func (c *Config) BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.Postgres.User,
		c.Postgres.Pass,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.DB,
	)
}
