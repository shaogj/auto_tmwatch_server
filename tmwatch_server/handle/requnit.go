package handle

import (
	log "202108FromBFLProj/ChainWatch_Project2023/auto_proj_test/logs"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
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
