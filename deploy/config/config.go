package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP      HTTPConfig      `mapstructure:"http"`
	Postgres  PostgresConfig  `mapstructure:"postgres"`
	CoinGecko CoinGeckoConfig `mapstructure:"coingecko"`
	Cron      CronConfig      `mapstructure:"cron"`
}

type HTTPConfig struct {
	Address string `mapstructure:"address"`
}
type PostgresConfig struct {
	URL string `mapstructure:"url"`
}
type CoinGeckoConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

type CronConfig struct {
	UpdateInterval string `mapstructure:"update_interval"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config: failed to read config")
	}
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "config: failed to unmarshal config")
	}
	return &config, nil
}
