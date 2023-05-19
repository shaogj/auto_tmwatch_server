package handle

import (
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/config"
	log "202108FromBFLProj/auto_tmwatch_server/tm_client_auto/logs"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//0505add,recover tm snapdata req
//二,,TM链的恢复
///home/dev-user/tmnode_snapdata0419.sh restoredata 20230423
//	r.POST("/sync_bsc", syncBsc)

type AccessToken struct {
	Token string `json:"token"`
}
type RecoverSnapDataRequest struct {
	AccessToken
	AutoIp       string `json:"auto_ip"`        //发出指令的服务IP
	Optype       string `json:"optype"`         //enum: restoredata:从节点指定的备份目录，拷贝数据到tm的node的目录
	SnapDataTime string `json:"snap_data_time"` //要恢复的定时备份的tm快照数据文件夹
}
type IPData struct {
	IPs    []string `json:"ips"`
	Type   string   `json:"type"` //tm or bsc or all
	Token  string   `json:"token"`
	Action string   `json:"action"` //add or del
}

type ExecResult struct {
	Stdout string
	Stderr string
	Cmderr string
}

func GetAppPath() string {
	path, err := os.Executable()
	if err != nil {
		fmt.Printf(err.Error())
	}
	dir := filepath.Dir(path)
	fmt.Printf(dir) // for example /home/user
	fmt.Println("cur gethapp run path", "dir", dir)
	return dir
}

// 0425
func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return "", nil //(cmd)
	} else {
		return runInLinux(cmd)
	}
}
func runInLinux(cmd string) (string, error) {
	fmt.Println("Running Linux cmd: %v", cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func ExecCmd(cmd *exec.Cmd) ExecResult {

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if stderr.Bytes() != nil {
		fmt.Println(stderr.String())
	}
	var errStr = ""
	if err != nil {
		fmt.Println(err.Error())
		errStr = err.Error()
		fmt.Println(errStr)
	}
	res := ExecResult{
		stdout.String(),
		stderr.String(),
		errStr,
	}
	return res

}

func sendTmSnapScriptCmd(optype, source, realsnaptimedir string) error {
	curPath := GetAppPath()
	timeToday := time.Now().Format("20060102_1504")
	fileTime := realsnaptimedir
	redirtmFile := fmt.Sprintf("%s//client_recoverruntm%s.log", curPath, timeToday)
	log.Logger.Infof("In sendScriptCmd(),check redirFile is:%s", redirtmFile)
	//./bscnode_snapdata_local.sh from 192.168.1.224 0525
	fileTime = "0525"
	trustConfigPath := fmt.Sprintf("%s//%s", curPath, "bscnode_snapdata_local.sh")
	log.Logger.Infof("to exec snapdataCheck file is :%s", trustConfigPath)
	curTotalCmd := fmt.Sprintf("%s %s %s %s > %s 2>&1", trustConfigPath, "from", "192.168.1.224", fileTime, redirtmFile)
	curRealShellCmd := fmt.Sprintf("%s//%s %s %s > %s 2>&1", curPath, "tmnode_snapdata0419.sh", optype, realsnaptimedir, redirtmFile)
	log.Logger.Infof("checking Online curRealShellCmd---> to exec tm cmd info is:%s", curRealShellCmd)

	//0511:/home/dev-user/tmnode_snapdata0419.sh restoredata 20230428

	log.Logger.Infof("start exec bscnode_snapdata.sh!,cmd is :%v:", curTotalCmd)
	pcmdres, err := RunCommand(curTotalCmd)
	log.Logger.Infof("after EXecCmd,get execResult info is :%v,err is:%v", pcmdres, err)
	return err
}

/*
curl --location --request POST '127.0.0.1:6667/sync_tm_snapdata' \
--header 'Content-Type: application/json' \
--data-raw '{"auto_ip":"192,135","optype":"restoredata","snap_data_time":"20230513","token":"4444"}'
*/
func SyncTmSnapData(c *gin.Context) {
	var syncDataRequest RecoverSnapDataRequest
	if err := c.BindJSON(&syncDataRequest); err != nil {
		log.Logger.Errorf("fun=SyncTmSnapData() requeset's Params Token is invalid! res err info=%v", err)
		return
	}

	log.Logger.Infof("fun=SyncTmSnapData()--receive sync tm snapdata request %+v", syncDataRequest)

	if syncDataRequest.Token != config.Conf.Service.AccessToken { //"4444"  //viper.GetString("token") {
		log.Logger.Errorf("fun=SyncTmSnapData() requeset's Params Token is invalid! req Token=%v,cfg's AccessToken is;%s ", syncDataRequest.Token, config.Conf.Service.AccessToken)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "async request SyncTmSnapData, StatusBadRequest token!"})
		return
	}
	if syncDataRequest.Optype != "restoredata" {
		log.Logger.Errorf("fun=SyncTmSnapData() requeset's Params Optype is invalid! req Optype=%v" + syncDataRequest.Optype)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "invalid params :Optype"})
		return
	}

	log.Logger.Info("async request sync data success! to handle request=%v", syncDataRequest)
	//0518add,note,//go SyncData?
	err := sendTmSnapScriptCmd(syncDataRequest.Optype, "localhost", syncDataRequest.SnapDataTime)
	if err != nil {
		log.Logger.Errorf("sync request sync data sendTmSnapScriptCmd() failed! cur handle request=%v", syncDataRequest)
		c.IndentedJSON(http.StatusOK, gin.H{"msg": "async request sync tm data failed!", "request resp err": err})
	}
	log.Logger.Info("cur handle request sync data sendTmSnapScriptCmd() finished! cur handle syncDataRequest.SnapDataTime=%v", syncDataRequest.SnapDataTime)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "async request sync tm data success!"})
}
func AddValidators(c *gin.Context) {
	log.Logger.Info("start AddValidators---PPP--AA", c.Request)
	var ipdata IPData

	if err := c.BindJSON(&ipdata); err != nil {
		return
	}
	log.Logger.Info("fun=AddValidators() bef--,request=%v", ipdata)
	//if ipdata.Token != config.Conf.Service
	if ipdata.Token != "4444" {
		log.Logger.Errorf("fun=AddValidators() requeset's Token is err,to break handle ,ipdata.Token=%v", ipdata)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "async request AddValidators, StatusBadRequest!"})
		return
	}
	log.Logger.Info("receive request addValidators success! to handle request=%v", ipdata)

	c.IndentedJSON(http.StatusOK, gin.H{"msg": "async request sync data success!"})

}
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		log.Logger.Info(c.Request.Header)
		log.Logger.Info(string(body))
		c.Next()
	}
}
