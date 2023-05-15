package handle

import (
	log "202108FromBFLProj/auto_tmwatch_server/tmwatch_server/logs"
	"strings"

	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net"
)

type IPData struct {
	IPs    []string `json:"ips"`
	Type   string   `json:"type"` //tm or bsc or all
	Token  string   `json:"token"`
	Action string   `json:"action"` //add or del
}

func AddValidators(c *gin.Context) {
	log.Logger.Info("start AddValidators---PPP--AA", c.Request)
	var ipdata IPData

	if err := c.BindJSON(&ipdata); err != nil {
		return
	}
	log.Info("fun=AddValidators() bef--,request=%v", ipdata)

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
