package main

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Welcome to go-template")

	// viper setup to read from config.toml
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("viper read config error")
	}

	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("started logger")
}
