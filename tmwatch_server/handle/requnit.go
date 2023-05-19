package handle

import (
	log "202108FromBFLProj/auto_tmwatch_server/tmwatch_server/logs"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"net"
)

type IPData struct {
	IPs    []string `json:"ips"`
	Type   string   `json:"type"` //tm or bsc or all
	Token  string   `json:"token"`
	Action string   `json:"action"` //add or del
}

func GetSnapDataTime() (datatimestr string) {
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
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GenTmErrMsg(errNodes []string, nodeType, msgPrefix string, clusterHeight int64) string {
	var msg DingErrMsg
	msg.MsgType = "text"
	for _, ip := range errNodes {
		var host HeightErrHost
		localHeight, _ := GetTmHeight(ip) //nodeType: tm
		host.ClusterHeight = clusterHeight
		host.LocalHeight = localHeight
		host.IP = ip
		host.Title = msgPrefix
		//0515ad
		msg.Text.AlarmLevelInfo = "tmwatch_SnapRecover&grade1"
		msg.Text.Content = append(msg.Text.Content, host)
	}
	content, _ := json.Marshal(msg)
	return string(content)
}

func GenOkMsg(okNodes []string, nodeType, msgPrefix string) string {
	var msg DingOkMsg
	msg.MsgType = "text"
	msg.Text.Content.Title = msgPrefix
	if len(okNodes) == 0 {
		msg.Text.Content.IPs = "所有节点"
	} else {
		msg.Text.Content.IPs = strings.Join(okNodes, ",")
	}
	content, _ := json.Marshal(msg)
	return string(content)
}
