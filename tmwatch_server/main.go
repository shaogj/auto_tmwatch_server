package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/config"
	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/handle"
	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/logs"
)

// 匹配新增的增量更新

func main() {

	config.LoadConf()
	log.LogInit(config.Conf.Service.LogLevel, config.Conf.Service.LogPath)
	handle.GLogger = log.Logger

	log.Logger.Infof("tm watch server--start!")
	log.Logger.Infof("get config.toml info ===>:%v:", config.Conf.Service)
	log.Logger.Infof("get config.toml for server TMClusterMonitor info ===>:%v:", config.Conf.TMClusterMonitor)
	log.Logger.Infof("get config.toml for TMClusterMonitor's TmMonitor is ===>:%v:", config.Conf.TmMonitor)
	log.Logger.Infof("get config.toml' Host'TM info ===>:%v:", config.Conf.TM)

	go handle.StartClusterStatusProc()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Logger.Infof("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //5*time.Second

	defer cancel()
	// catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		log.Logger.Infoln("timeout of 2 seconds.")
	}
	log.Logger.Info("auto_tmwatch_server exiting")

}
