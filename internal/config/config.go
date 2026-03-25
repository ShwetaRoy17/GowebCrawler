package config


// import (
// 	"github.com/spf13/viper"
// )

type Config struct {
	seedUrls []string `mapstructure:"seed_urls"`
	MaxDepth int `mapstructure:"max_depth"`
	Concurrency int `mapstructure:"concurrency"`
}