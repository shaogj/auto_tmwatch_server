package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/config"
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/handle"
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/logs"
)

type AccessToken struct {
	Token string `json:"token"`
}

type CompressRequest struct {
	AccessToken
	FileTime string `json:"fileTime"`
	AutoIp   string `json:"autoIp"`
}

func main() {
	curFormtTime := time.Now().Format("2006.01.02_15_04")
	fmt.Printf("cur curFormtTime is:%s\n", curFormtTime)

	//sendScriptCmd()
	config.LoadConf()
	log.LogInit(config.Conf.Service.LogLevel, config.Conf.Service.LogPath)

	log.Logger.Infof("tm watch server--start!")
	log.Logger.Infof("get config.toml info ===>:%v:", config.Conf.Service)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//r.Use(handle.RequestLoggerMiddleware())
	r.POST("/add_validators", handle.AddValidators)

	r.POST("/sync_tm_snapdata", handle.SyncTmSnapData)
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

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Logger.Infof("Shutdown Server ...")

}
