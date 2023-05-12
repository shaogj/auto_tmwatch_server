package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/config"
	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/handle"
	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/logs"
	"github.com/gin-gonic/gin"
)

// 匹配新增的增量更新

func main() {

	config.LoadConf()
	log.LogInit(config.Conf.Service.LogLevel, config.Conf.Service.LogPath)
	handle.GLogger = log.Logger

	log.Logger.Infof("tm watch server--start!")
	log.Logger.Infof("get config.toml info ===>:%v:", config.Conf.Service)
	log.Logger.Infof("get config.toml' TMClusterMonitor info ===>:%v:", config.Conf.TMClusterMonitor)
	log.Logger.Infof("get config.toml' Host'TM info ===>:%v:", config.Conf.TM)
	//level := config.Conf.Service.LogLevel
	/*	go TMWatch()
		go MonitorBscCluster()
	*/
	log.Logger.Info("to run StartClusterStatusProc()")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(handle.RequestLoggerMiddleware())
	r.POST("/add_validators", handle.AddValidators)
	port := config.Conf.Service.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//0426ing
	go handle.StartClusterStatusProc()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//0427doing
	//skip---config.SaveConf(*config.Conf)

	log.Logger.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //5*time.Second
	//2023.0401

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Logger.Infoln("timeout of 5 seconds.")
	}
	log.Logger.Info("Server exiting")

}
