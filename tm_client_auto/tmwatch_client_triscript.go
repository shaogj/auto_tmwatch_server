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

func getSnapDataTime() (datatimestr string) {
	var getdatatimestr string
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	timestr := addTime.Format("2006-01-02")
	addTime2hour := time.Date(t.Year(), t.Month(), t.Day(), 2, 0, 0, 0, t.Location())
	fmt.Println("today 2hour is:%s,timestemp is:%d ", timestr, addTime2hour.Unix()) //2022-04-15
	tyesday := time.Now().AddDate(0, 0, -1)
	addTimeyesday2hour := time.Date(tyesday.Year(), tyesday.Month(), tyesday.Day(), 2, 0, 0, 0, t.Location())
	timeyesdaystr := addTimeyesday2hour.Format("2006-01-02")
	fmt.Println("yestorday 2hour is:%s,timestemp is:%d ", timeyesdaystr, addTimeyesday2hour.Unix()) //2022-04-15
	curtime := time.Now().Unix()
	if curtime > addTime2hour.Unix() {
		getdatatimestr = addTime.Format("20060102")
	} else {
		addTimeyesday0hour := time.Date(tyesday.Year(), tyesday.Month(), tyesday.Day(), 0, 0, 0, 0, t.Location())
		getdatatimestr = addTimeyesday0hour.Format("20060102")
	}
	return getdatatimestr
}
func main() {
	curFormatDayTime := time.Now().Format("2006.01.02_15_04")
	curFormatTime := time.Now().Format("2006.01.02_15_04")
	fmt.Printf("cur curFormatDayTime is:%s,curFormatTime is:%s\n", curFormatDayTime, curFormatTime)
	getdatatimestr := getSnapDataTime()
	fmt.Printf("cur nowtime is is::%d,getSnapDataTime() is:%s\n", time.Now().Unix(), getdatatimestr)

	return
	//0519testing
	//sendScriptCmd()
	config.LoadConf()
	log.LogInit(config.Conf.Service.LogLevel, config.Conf.Service.LogPath)

	log.Logger.Infof("tm watch server--start!")
	log.Logger.Infof("get config.toml of client info ===>:%v:", config.Conf.Service)

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
