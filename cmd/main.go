package main

import (
	"context"
	"flag"
	"github.com/quanxiang-cloud/cabin/logger"

	"os"
	"os/signal"
	"syscall"

	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/goalie/api/restful"
	"github.com/quanxiang-cloud/goalie/pkg/config"
)

var (
	configPath = flag.String("config", "../configs/config.yml", "-config 配置文件地址")
)

func main() {
	flag.Parse()
	log := logger.Logger
	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}

	db, err := mysql2.New(conf.Mysql, log)
	if err != nil {
		panic(err)
	}
	// 启动路由
	ctx := context.Background()
	router, err := restful.NewRouter(ctx, conf, log, db)
	if err != nil {
		panic(err)
	}
	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			router.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
