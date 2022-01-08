package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"tp-db-project/configs"
	"tp-db-project/internal/pkg/server"
	"tp-db-project/internal/pkg/utilits"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/configs.toml", "path to config file")
}

func main() {
	fmt.Println("start main")
	flag.Parse()
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := utilits.NewLogger(config)
	//defer func(closer func() error, log *logrus.Logger) {
	//	err := closer()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}(logger)

	db, closeDbResource := utilits.NewPostgresConnection(&config.ServerRepository)

	defer func(closer func(), log *logrus.Logger) {
		closer()
	}(closeDbResource, logger)

	serv := server.NewServer(config,
		utilits.ExpectedConnections{
			SqlConnection: db,
		},
		logger,
	)
	if err = serv.Start(); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Server was stopped")

}
