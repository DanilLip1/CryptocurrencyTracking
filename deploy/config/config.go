package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddress        string `mapstructure:"http_address"`
	PostgresURL        string `mapstructure:"postgres_url"`
	CoinGeckoApiKey    string `mapstructure:"coingecko_api_key"`
	CoinGeckoBaseURL   string `mapstructure:"coingecko_base_url"`
	CronUpdateInterval string `mapstructure:"cron_update_interval"`
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
