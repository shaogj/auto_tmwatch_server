package unit_req

import (
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/config"
	"encoding/json"
	"fmt"
	"github.com/mkideal/log"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetChainHeight() (int64, error) {
	return 0, nil
}

func post(url string, payload *strings.Reader) ([]byte, error) {

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetBSCChainHeight() (int64, error) {
	type Result struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  string `json:"result"`
	}
	var res Result
	//url := "http://127.0.0.1:8545"
	url := "http://106.3.133.180:8545"
	//url := "http://192.168.1.224:8545"

	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber", "id":1}`)

	b, err := post(url, payload)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return 0, err
	}

	res.Result = strings.TrimPrefix(res.Result, "0x")

	if res.Result == "" {
		res.Result = "0"
	}

	h, err := strconv.ParseInt(res.Result, 16, 64)
	if err != nil {
		return 0, err
	}
	return h, nil

}

//0104add

func get(url string) ([]byte, error) {
	// url := "http://106.3.133.179:46657/tri_block_info?height=104360"
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetTmHeight(url string) (int64, error) {
	type jsonResult struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  struct {
			Response struct {
				Data             string `json:"data"`
				Version          string `json:"version"`
				AppVersion       string `json:"app_version"`
				LastBlockHeight  string `json:"last_block_height"`
				LastBlockAppHash string `json:"last_block_app_hash"`
			} `json:"response"`
		} `json:"result"`
	}
	var jsonRet jsonResult
	ret, err := get(url)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(ret, &jsonRet)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(jsonRet.Result.Response.LastBlockHeight, 10, 64)
}

func MonitorBSCBlockIncrease() {
	var oldH int64
	// wait chain start
	//time.Sleep(10 * time.Second)
	errTimes := 0
	maxErr := 5
	//tmp check 0107
	interval := 15 //7 //15 //60
	for {
		fmt.Printf("111--start\n")
		newH, err := GetBSCChainHeight()
		fmt.Printf(fmt.Sprintf("cur times ,req GetBSCChainHeight():get newH val is: %d,err is :%v\r\n", newH, err))

		if err == nil {
			if !(newH > oldH) {
				//0105update
				fmt.Printf(fmt.Sprintf("MonitorBlockIncrease: block height dont inscrase,interval %d\n", interval))
				panic(fmt.Sprintf("MonitorBlockIncrease: block height dont inscrase,interval %ds\n", interval))
			} else {
				oldH = newH
			}
		} else {
			fmt.Printf(fmt.Sprintf("get height err: %s", err.Error()))

			if errTimes > maxErr {
				fmt.Printf(fmt.Sprintf("MonitorBlockIncrease: err times  greater than %d\n", maxErr))
				panic(fmt.Sprintf("MonitorBlockIncrease: err times  greater than %d", maxErr))
			}
			errTimes += 1
		}
		fmt.Printf("111--end cur times\n")

		time.Sleep(time.Duration(interval) * time.Second)

	}
}

// 0112add,[1000--10000]
func CheckRandnum() {
	inum := 0
	for inum < 3 {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(4000) + 1000
		fmt.Println(r)
		inum++

	}
}

func MonitorTMBlockIncrease(curruncfg *config.ConfigMonitorInfomationTM, localip string) {
	var oldH int64
	errTimes := 0
	maxErr := 5
	interval := curruncfg.RequestTmInterval //60 //15 //60
	tmnodeip := "192.168.1.225"
	//tmnodeip = curruncfg.ReqRpcUrl
	tmnodeip = localip
	reqTMNodeUrl1 := fmt.Sprintf("http://%s:46657/tri_abci_info?", tmnodeip)
	itimes := 0
	watchprocname := "tendermint"
	for {
		itimes++
		log.Info("222--start check tmnode height,reqTMNodeUrl1 is:%s,req RequestTmInterval num is :%d,RequestTmInterval is :%d\n", reqTMNodeUrl1, interval, curruncfg.RequestTmInterval)
		newH, err := GetTmHeight(reqTMNodeUrl1)
		log.Info(fmt.Sprintf("curreq times ,GetTMChainHeight()' get newH val is: %d,errTimes is:%d,err is :%v\r\n", newH, errTimes, err))
		if err == nil {
			//0113,add confi
			if !(newH > oldH) { //|| itimes > 3 {
				log.Error(fmt.Sprintf("MonitorTMBlockIncrease: block height dont increase,interval %d,cur newH is:%d, oldH is:%d\n", interval, newH, oldH))
				//0217enhance,watch n times no inscrease:
				errTimes++
				if errTimes < maxErr {
					log.Error("cur checking errTimes is:%d,to sleep RequestTmInterval is:%d", errTimes, curruncfg.RequestTmInterval)
					time.Sleep(time.Duration(curruncfg.RequestTmInterval) * time.Second)
					continue
				}
				log.Error("cur checking errTimes is more then :%d,to exec KillProcessByName().to restart the tmnode!", maxErr)
				log.Info("step1:checking===start KillProcessByName is:%v\n", watchprocname)
				error := KillProcessByName(watchprocname)
				if error != nil {
					log.Error("checking===exec KillProcessByName() failed! kill watchprocname is:%v,err is:%v\n", watchprocname, error)
					continue
				} else {
					//exist no the pid is,also KillProcessByName return nil
					log.Info("checking===exec KillProcessByName() succ!  watchprocname : %v is killed,err is:%v\n", watchprocname, error)
				}
				//2023.0307
				if curruncfg.StartMode == 0 {
					log.Info("step2:to exec StopTendermint(),watchprocname is:%v\n", watchprocname)
					pid, err := StopTendermint("tendermint")
					log.Info("checking===StopTendermint() is done. get is pid is:%v,err is:%v\n", pid, err)

					log.Info("step3:to exec StartTendermint(),to after sleep OperatorSysInterval second is:%d", 5)
					//curruncfg.OperatorSysInterval
					time.Sleep(time.Duration(4) * time.Second)
					pid, err = StartTendermint("tendermint")
					//fmt.Printf
					log.Info("checking===StartTendermint() is done.get is pid is:%v,err is:%v\n", pid, err)
				} else {
					//StartMode =1,cmd run server mode
					go func() {
						time.Sleep(5 * time.Second)
						log.Info("step2:step2:to exec start tendermint node In Mode2.to run cmd, watchprocname is:%v\n", watchprocname)
						StartTMServer()
					}()
				}
				log.Info("for start Tendermint Server.cur delay to sleep OperatorSysInterval is:%d", curruncfg.OperatorSysInterval)
				//0220,add getpubkey(),to delay is 40s;
				time.Sleep(time.Duration(curruncfg.OperatorSysInterval) * time.Second)

			} else {
				oldH = newH
			}
			//0217enhance
			errTimes = 0

		} else {
			errTimes += 1
			log.Error(fmt.Sprintf("cur req GetTmHeight() get height errTimes is:%d, err: %v\n", errTimes, err.Error()))
			if errTimes > maxErr {
				log.Info(fmt.Sprintf("MonitorTMBlockIncrease: err times  greater than %d,cur newH is:%d\n", maxErr, newH))
				//panic(fmt.Sprintf("MonitorTMBlockIncrease: err times  greater than %d", maxErr))
				log.Info("step1:In expeed errTimes!,checking===start KillProcessByName is:%v\n", watchprocname)
				error := KillProcessByName(watchprocname)
				if error != nil {
					log.Error("checking===exec KillProcessByName() failed! kill watchprocname is:%v,err is:%v\n", watchprocname, error)
				}
				//2023.0307new add
				if curruncfg.StartMode == 0 {
					log.Info("step2:to exec StopTendermint(),watchprocname is:%v\n", watchprocname)
					pid, err := StopTendermint("tendermint")
					log.Info("checking===StopTendermint() is done. get is pid is:%v,err is:%v\n", pid, err)
					time.Sleep(time.Duration(4) * time.Second)

					log.Info("step3:In expeed errTimes!to exec StartTendermint()")
					pid, err = StartTendermint("tendermint")
					if err != nil {
						log.Error("checking===exec StartTendermint() failed! kill watchprocname is:%v,err is:%v\n", watchprocname, error)
					} else {
						log.Info("checking===StartTendermint() is done.get is pid is:%v,err is:%v\n", pid, err)
					}
				} else {
					//StartMode =1,cmd run server mode
					go func() {
						//bj env
						time.Sleep(5 * time.Second)
						log.Info("step2:to exec start tendermint node In Mode2. to run cmd, watchprocname is:%v\n", watchprocname)
						StartTMServer()
					}()
				}
				log.Info("for start Tendermint Server.cur delay to sleep OperatorSysInterval is:%d", curruncfg.OperatorSysInterval)
				//0220,add getpubkey(),to delay is 40s;
				time.Sleep(time.Duration(curruncfg.OperatorSysInterval) * time.Second)

				errTimes = 0
			}

		}
		log.Info("222--end check tmnode height,cur tm times.to sleep RequestTmInterval is:%d\n", curruncfg.RequestTmInterval)
		time.Sleep(time.Duration(curruncfg.RequestTmInterval) * time.Second)

		CheckRandnum()
		//time.Sleep(time.Duration(interval) * time.Second)

	}
}
