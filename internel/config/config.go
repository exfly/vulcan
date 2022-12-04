package config

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

//go:generate gomodifytags -file config.go -all -add-tags mapstructure -w -override -transform snakecase -quiet

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment          string        `mapstructure:"environment"`
	DBDriver             string        `mapstructure:"db_driver"`
	DBSource             string        `mapstructure:"db_source"`
	RedisAddress         string        `mapstructure:"redis_address"`
	HTTPServerAddress    string        `mapstructure:"http_server_address"`
	GRPCServerAddress    string        `mapstructure:"grpc_server_address"`
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

type Path string

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path Path) (config *Config, err error) {
	viper.AddConfigPath(string(path))
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	config = &Config{}
	err = viper.Unmarshal(config)
	spew.Dump(config)
	return
}
