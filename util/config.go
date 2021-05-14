package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBDataSource        string        `mapstructure:"DB_DATA_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	WXAppid             string        `mapstructure:"WX_APPID"`
	WXSecret            string        `mapstructure:"WX_SECRET"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func NewConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app-dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
