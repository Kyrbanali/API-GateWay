package config

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
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

func Load() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("read config.yaml", slog.Any("err", err))
	}

	if err := godotenv.Load(); err != nil {
		slog.Warn(".env not found", slog.Any("err", err))
	}

	viper.AutomaticEnv()

	cfg := &Config{
		Postgres: PostgresConfig{
			Host: viper.GetString("POSTGRES_HOST"),
			User: viper.GetString("POSTGRES_USER"),
			Pass: viper.GetString("POSTGRES_PASS"),
			DB:   viper.GetString("POSTGRES_DB"),
			Port: viper.GetString("POSTGRES_PORT"),
		},
		App: AppConfig{
			Port: viper.GetString("API_PORT"),
		},
	}

	return cfg, nil
}
