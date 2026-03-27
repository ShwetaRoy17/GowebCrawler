package config


// import (
// 	"github.com/spf13/viper"
// )

type Config struct {
	SeedUrl string `mapstructure:"seedurl"`
	MaxDepth int `mapstructure:"max_depth"`
	Concurrency int `mapstructure:"concurrency"`
}

