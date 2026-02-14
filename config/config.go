package config

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/Kyrbanali/API-GateWay/internal/cache"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
	"github.com/Kyrbanali/API-GateWay/internal/worker"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port string
	}

	Postgres storage.Config
	Cache    cache.Config
	Workers  worker.Config

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
