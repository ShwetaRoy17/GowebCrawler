package config

import(
	 "github.com/spf13/viper"
	 "fmt"
)


type Config struct {
	SeedUrl string `mapstructure:"seedurl"`
	MaxDepth int `mapstructure:"max_depth"`
	Concurrency int `mapstructure:"concurrency"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("max_depth",3)
	viper.SetDefault("concurrency",10)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal : %w",err)
	}
	return cfg, nil

}

