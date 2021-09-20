package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gbaeke/go-template/pkg/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	fmt.Printf("Welcome to go-template\n\n")

	// flags
	f := pflag.NewFlagSet("api", pflag.ContinueOnError)
	// web server port
	f.Int("port", 8080, "HTTP Port")
	// message to print on /
	f.String("welcome", "hello", "Welcome Message")
	// turn request logging on or off
	f.Bool("log", false, "Turn HTTP logging on or off")
	// graceful server shutdown timeout; use --timeout 15s etc...
	f.Duration("timeout", time.Second*15, "Server graceful shutdown timeout")

	// parse flags and display help message on -help
	err := f.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error %s\n", err.Error())
	}

	// bind flags to viper
	viper.BindPFlags(f)

	//also read from config.toml - flags come first
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "viper read config error %s\n", err.Error())
	}

	// init zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("started logger")

	// unmarshal config in api.Config struct
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		sugar.Panic("could not unmarshal config", zap.Error(err))
	}

	// log the config values - just an example of logging
	sugar.Infow("config values",
		zap.String("welcome", srvCfg.Welcome),
		zap.Int("port", srvCfg.Port),
		zap.Bool("log", srvCfg.Log),
		zap.Duration("timeout", srvCfg.Timeout),
	)

	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, sugar)
	srv.StartServer()

}
