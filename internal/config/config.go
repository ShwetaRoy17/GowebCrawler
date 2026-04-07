package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	SeedUrl     string `mapstructure:"seedurl"`
	MaxDepth    int    `mapstructure:"max_depth"`
	Concurrency int    `mapstructure:"concurrency"`
	UserAgent   string `mapstructure:"user_agent"`
	RateLimit   int    `mapstructure:"rate_limit"`
	Burst       int    `mapstructure:"burst"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("max_depth", 3)
	viper.SetDefault("concurrency", 10)
	viper.SetDefault("user_agent", "goWebC/1.0")
	viper.SetDefault("rate_limit", 10)
	viper.SetDefault("burst", 6)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal : %w", err)
	}
	return cfg, nil

}
