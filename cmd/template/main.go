package main

import (
	"fmt"

	"github.com/gbaeke/go-template/pkg/api"
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

	// retrieve config
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		sugar.Panic("could not unmarshal config", zap.Error(err))
	}

	// log the config values
	sugar.Infow("config values",
		zap.String("welcome", srvCfg.Welcome),
		zap.Int("port", srvCfg.Port),
	)

	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, sugar)
	srv.StartServer()

}
