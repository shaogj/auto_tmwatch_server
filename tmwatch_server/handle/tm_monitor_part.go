package handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
func GetClusterHStatus(nodeType string) (ClusterHStatus, error) {
	var cluster ClusterHStatus

	ips := GetIps(nodeType)
	//0426,checking, from fmt.Println
	fmt.Printf("to check tm ip list len is:%d, info is:%s\n", len(ips), ips)
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
				fmt.Printf("cur GetTmHeight() res error! ,cur node ip:%s, height is:%d:%v\n", ip, h, err)
				results <- chResult
				return
			}
			chResult.H.Height = h
			fmt.Printf("get GetTmHeight() res info : serion id :%d,tmHeight: %s:%d\n", ino, chResult.H.Ip, chResult.H.Height)
			results <- chResult
		}(ino, tmIp)
	}

	for i := 0; i < len(ips); i++ {
		result := <-results
		fmt.Println("get cur chan result infois :%v", result)

		cluster.Nodes = append(cluster.Nodes, result.H)

	}
	close(results)
	return cluster, nil
}
