package main

import (
	"202108FromBFLProj/ChainWatch_Project2023/auto_proj_test/config"
	uselog "202108FromBFLProj/ChainWatch_Project2023/auto_proj_test/logs"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"202108FromBFLProj/ChainWatch_Project2023/auto_proj_test/handle"
	"github.com/gin-gonic/gin"
)

// 匹配新增的增量更新

func main() {

	// go ConfigWatch()
	uselog.LogInit(config.Conf.Service.LogLevel, config.Conf.Service.LogPath)
	//level := config.Conf.Service.LogLevel

	/*	go TMWatch()
		go BscWatch()
		go MonitorBscCluster()
	*/
	uselog.Logger.Info("start NodeUpdateServer")
	//fmt.Printf("Version: %s\n", Version)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(handle.RequestLoggerMiddleware())
	r.POST("/add_validators", handle.AddValidators)
	//r.GET("/get_validators", GetValidators)
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

	go func() {
		times := 0
		errtimes := 0
		retry := 3
		var lastMaxHeight int64
		var newMaxHeight int64
		for times < retry {
			tms, err := handle.GetClusterHStatus("tm")
			fmt.Printf("after check all tm's IPlist,by GetClusterHStatus() getinfo: %v,err is:%v \n", tms, err)
			//0504checking
			time.Sleep(time.Duration(7) * time.Second)
			newMaxHeight = handle.GetMaxH(tms)
			if newMaxHeight > lastMaxHeight {
				lastMaxHeight = newMaxHeight
				fmt.Printf("after GetClusterHStatus(), check tmchain height is increase! get newMaxHeight is :%d,lastMaxHeight is:%d,\n", lastMaxHeight, newMaxHeight)
			} else {
				errtimes++
				fmt.Printf("after GetClusterHStatus() ,check tmchain height is increase no ! errtimes is :%d,get newMaxHeight is :%d,lastMaxHeight is:%d,\n", errtimes, lastMaxHeight, newMaxHeight)
			}
			times++
		}

	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//0427doing
	//skip---config.SaveConf(*config.Conf)

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//2023.0401
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
