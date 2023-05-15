package handle

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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
