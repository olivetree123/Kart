package config

import (
	"github.com/spf13/viper"
)

// Config 配置
var Config *viper.Viper

func init() {
	Config = viper.New()
	Config.SetConfigName("kart")
	Config.AddConfigPath("./")
	err := Config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
