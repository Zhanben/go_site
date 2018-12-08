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

// @title Go-site Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1
// @BasePath ""
func main() {

	//Read config file
	err := config.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read config file: %s \n", err))
	}

	//Init log
	log.InitLog()

	//Init db connection
	dbConn, err := db.InitDbConn()
	if err != nil {
		fmt.Printf("connect db error:%s", err)
		panic("connect db error!")
	}
	log.Logger.Info("Db init successful!")


	exit := make(chan os.Signal,10) //初始化一个channel
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM) //notify方法用来监听收到的信号
	//start http sever
	go func() {
		startServer(dbConn)
	}()
	sig := <-exit
	log.Logger.Info("main function return, %s", sig.String())


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
