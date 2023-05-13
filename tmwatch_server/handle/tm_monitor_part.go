package handle

import (
	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/config"
	"202108FromBFLProj/auto_tmwatch_server/tmwatch_server/logs"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var GLogger *zap.SugaredLogger

// 节点块高
type NodeH struct {
	Height int64
	Ip     string
}

// 集群的块高状态
type ClusterHStatus struct {
	Nodes []NodeH
}

// 获取节点块高,channel返回结果
type ChResult struct {
	H   NodeH
	Err error
}

// 落后节点map
type LagNodes map[string]int

// 高度错误节点
type HeightErrHost struct {
	Title         string
	IP            string
	LocalHeight   int64
	ClusterHeight int64
}

type HeightOkHosts struct {
	Title string
	IPs   string
}

// 钉钉异常通知text
type ErrText struct {
	Content []HeightErrHost `json:"content"`
}

var (
	ErrBscAddr = errors.New("invaild bsc addr")
	ErrTmAddr  = errors.New("invaild tm addr")
)

func post(url string, payload *strings.Reader) ([]byte, error) {

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

//0426add，to send tmwatch cmdscript run
// 发送tmwatch恢复快照功能

func sendMsg(url, nodeType, content string) (getret string, err error) {
	fmt.Printf("send", nodeType, "Msg:", content, "url:", url)
	payload := strings.NewReader(content)
	ret, err := post(url, payload)
	return string(ret), err
	//var dingRet DingResp
	/*
		json.Unmarshal(ret, &dingRet)
		if dingRet.ErrCode != 0 {
			fmt.Println("钉钉调用错误: %v", dingRet)
		}
	*/
}
func get(url string) ([]byte, error) {
	// url := "http://106.3.133.179:46657/tri_block_info?height=104360"

	client := &http.Client{}
	client.Timeout = time.Second * 60
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("%s", err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("%s", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetTmHeight(ip string) (int64, error) {

	var res TmHResponse
	url := "http://" + ip + ":46657/tri_abci_info"
	r, err := get(url)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(r, &res)
	if err != nil {
		fmt.Println("%s", r)
		fmt.Println(ip)
		return 0, err
	}
	if res.Result.Response.LastBlockHeight == "" {
		res.Result.Response.LastBlockHeight = "0"
	}
	h, err := strconv.ParseInt(res.Result.Response.LastBlockHeight, 10, 64)
	if err != nil {
		return 0, err
	}
	return h, nil

}

func GetMaxH(tms ClusterHStatus) int64 {
	var max int64
	for _, host := range tms.Nodes {
		if host.Height > max {
			max = host.Height
		}
	}
	return max
}
func GetIps(nodeType string) []string {
	return []string{"148.153.184.138", "164.52.93.91", "148.153.184.133"}
}

func GetBjIps(nodeType string) []string {
	return []string{"106.3.133.178", "106.3.133.179", "106.3.133.180", "210.73.218.171", "210.73.218.172"}
}
func GetClusterHStatus(nodeType string, tmnodelist []string) (ClusterHStatus, error) {
	var cluster ClusterHStatus

	//0512doing,bjenv
	//ips := GetBjIps(nodeType)
	ips := tmnodelist
	//0426,checking, from fmt.Println
	//fmt.Printf
	GLogger.Infof("to check tm GetBjIps()' ip list len is:%d, info is:%s\n", len(ips), ips)
	results := make(chan ChResult, len(ips))
	for ino, tmIp := range ips {
		go func(ino int, ip string) {
			var chResult ChResult
			chResult.H.Ip = ip

			h, err := GetTmHeight(ip)
			//0426doing
			if err != nil {
				fmt.Println(err)
				fmt.Println(ip)
				chResult.Err = err
				//0512,to mark err ip
				GLogger.Errorf("cur GetTmHeight() res error! ,cur node ip:%s, height is:%d:%v\n", ip, h, err)
				results <- chResult
				return
			}
			chResult.H.Height = h
			GLogger.Infof("get GetTmHeight() res info : serion id :%d,tmHeight: %s:%d\n", ino, chResult.H.Ip, chResult.H.Height)
			results <- chResult
		}(ino, tmIp)
	}

	for i := 0; i < len(ips); i++ {
		result := <-results
		GLogger.Infof("get cur chan result infois :%v", result)

		cluster.Nodes = append(cluster.Nodes, result.H)

	}
	close(results)
	return cluster, nil
}

func StartClusterStatusProc() {
	//0513test
	/*good invoke
	url := "http://127.0.0.1:6667/sync_tm_snapdata"
	fileTime := "0512datanoon"
	fmt.Printf("cur check SendTmSnapRecoverRequest()!")
	//SendCompressBscRequest(url, fileTime)
	SendTmSnapRecoverRequest(url, fileTime)
	return
	
	*/
	times := 0
	errtimes := 0
	retry := 3
	var lastMaxHeight int64
	var newMaxHeight int64
	var tmnodelist []string
	for _, host := range config.Conf.TM {
		tmnodelist = append(tmnodelist, host.Ip)
	}
	//0512doing,,//log.Logger
	GLogger.Infof("cur GetClusterHStatus(), get tm's IPlist: %v\n", tmnodelist)
	//fmt.Printf("cur GetClusterHStatus(), get tm's IPlist: %v\n", tmnodelist)

	//return
	for times < retry {
		tms, err := GetClusterHStatus("tm", tmnodelist)
		GLogger.Infof("after check all tm's IPlist,by GetClusterHStatus() getinfo: %v,err is:%v \n", tms, err)
		//0504checking
		time.Sleep(time.Duration(7) * time.Second)
		newMaxHeight = GetMaxH(tms)
		if newMaxHeight > lastMaxHeight {
			//fmt.Printf
			GLogger.Infof("after GetClusterHStatus(), check tmchain height is increase! get newMaxHeight is :%d,lastMaxHeight is:%d,\n", newMaxHeight, lastMaxHeight)
			lastMaxHeight = newMaxHeight
		} else {
			errtimes++
			GLogger.Errorf("after GetClusterHStatus() ,check tmchain height is increase no ! errtimes is :%d,get newMaxHeight is :%d,lastMaxHeight is:%d,\n", errtimes, newMaxHeight, lastMaxHeight)
			//0512add
			if errtimes > 2 {
				//POST '127.0.0.1:6667/sync_tm_snapdata'
				url := "http://127.0.0.1:6667/sync_tm_snapdata"
				fileTime := "0512datanoon"
				GLogger.Infof("cur check tmchain is increase no errtimes is :%d,newMaxHeight is :%d,to invoke SendCompressBscRequest()\n", errtimes, newMaxHeight)
				//SendCompressBscRequest(url, fileTime)
				SendTmSnapRecoverRequest(url, fileTime)
			}
		}
		times++
	}

	//}()

}

func SendCompressBscRequest(url, fileTime string) {
	autoIp := GetOutboundIP().String()

	payload := strings.NewReader(fmt.Sprintf(`{"token":"%s","autoIp":"%s", "fileTime":"%s"}`, config.Conf.TMClusterMonitor.AccessToken, autoIp, fileTime))
	log.Logger.Infof("send snaprecover for tm' request :%+v", payload)
	post(url, payload)

}

func SendTmSnapRecoverRequest(url, fileTime string) {
	autoIp := GetOutboundIP().String()
	//{"auto_ip":"192,135","optype":"restoredata","snap_data_time":"20230513","token":"4444"}
	payload := strings.NewReader(fmt.Sprintf(`{"auto_ip":"%s","optype":"%s","snap_data_time":"%s", "token":"%s"}`, autoIp, "restoredata", fileTime, config.Conf.TMClusterMonitor.AccessToken))
	log.Logger.Infof("send tmnode snaprecover for tm' request :%+v", payload)
	post(url, payload)

}
