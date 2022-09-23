package main

import (
	"fmt"
	"go-redis/config"
	"go-redis/lib/logger"
	"go-redis/resp/handler"
	"go-redis/tcp"
	"os"
)

const configFile string = "redis.conf"

var defaultProperites = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExist(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	//logger.Setup(&logger.Settings{
	//	Path:       "logs",
	//	Name:       "go-redis",
	//	Ext:        "log",
	//	TimeFormat: "2022-08-28",
	//})

	if fileExist(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperites
	}
	//err := tcp.ListenAndServeWithSignal(&tcp.Config{
	//	Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	//}, tcp.NewEchoHandler())
	//if err != nil {
	//	logger.Fatal(err)
	//}

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, handler.MakeHandler())
	if err != nil {
		logger.Fatal(err)
	}
}
