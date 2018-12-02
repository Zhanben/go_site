package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhanben/go_site/app"
	"github.com/zhanben/go_site/tool/config"
	"github.com/zhanben/go_site/tool/db"
	"github.com/zhanben/go_site/tool/log"

	"go.uber.org/zap"
)

func main() {
	//Read config file
	err := config.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read config file: %s \n", err))
	}

	//Init log
	log.InitLog()

	//Init db connection
	db, err := db.InitDbConn()
	if err != nil {
		panic("connect db error!")
	}
	log.Logger.Info("Db init successful!")

	//start http sever
	startServer(db)
}

func startServer(dbConn *db.DbConn) error {
	exit := make(chan os.Signal, 10)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	s, err := app.NewServer(dbConn, log.Logger)
	if err != nil {
		log.Logger.Panic("init http server failed.")
	}
	err = s.HttpServer.ListenAndServe()
	if err != nil {
		log.Logger.Panic("start http server error:%s", zap.Error(err))
	}
	return nil
}
