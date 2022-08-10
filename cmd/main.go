package main

import (
	"fmt"
	"os"
	"os/signal"

	server "github.com/DarkSoul94/money-processing-service/cmd/httpserver"
	"github.com/DarkSoul94/money-processing-service/config"
	"github.com/DarkSoul94/money-processing-service/pkg/logger"
)

func main() {
	conf := config.InitConfig()
	logger.InitLogger(conf)

	apphttp := server.NewApp(conf)
	go apphttp.Run(conf)

	fmt.Println(
		fmt.Sprintf(
			"Service %s is running",
			conf.AppName,
		),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	apphttp.Stop()

	fmt.Println(
		fmt.Sprintf(
			"Service %s is stopped",
			conf.AppName,
		),
	)
}
