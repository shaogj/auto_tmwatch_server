package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
)

var Conf = &Config{}

type Host struct {
	Ip string `toml:"ip" comment:"ip地址"`
	//IsNew bool   `toml:"isNew" comment:"用于确定是否为新加ip,默认为false,已弃用"`
}

type TmConf struct {
	IP string `toml:"ip"`
}

type BscConf struct {
	Token string `toml:"token" comment:"bsc链server token"`
}

type Monitor struct {
	DingUrl      string `toml:"ding-url" comment:"钉钉机器人接口"`
	OkPrefixKey  string `toml:"ok-prefix-key" comment:"钉钉安全设置,节点正常消息前缀"`
	ErrPrefixKey string `toml:"err-prefix-key" comment:"钉钉安全设置,节点异常消息前缀"`
	Interval     int    `toml:"interval" comment:"获取块高时间间隔,单位为秒"`
	RetryTimes   int    `toml:"retry-times" comment:"判断区块落后所需次数"`
	//AbnormalHosts []string `toml:"abnormal-hosts" comment:"落后区块节点ip列表"`
}

type ServiceConf struct {
	Port     int    `toml:"port" comment:"服务端口默认为6667"`
	LogLevel string `toml:"log-level" comment:"日志等级int类型默认为info; debug || info || warn || error"`
	LogPath  string `toml:"log-path" comment:"日志路径,默认为update.log"`
}

type ClusterMonitor struct {
	ClusterInterval int    `toml:"cluster-interval" comment:"查看集群状态时间间隔,默认10分钟"`
	NodeInterval    int    `toml:"node-interval" comment:"bsc节点重启时间间隔,默认2分钟"`
	MonitorRpc      int    `toml:"monitor-rpc" comment:"用于重启bsc的rpc端口"`
	AccessToken     string `toml:"access-token" comment:"访问tmwatch-client的请求token"`
}
type Config struct {
	TM               []Host         `toml:"tm" comment:"tm链节点ip列表"`
	TmMonitor        Monitor        `toml:"tm-monitor" comment:"tm落后节点监视器"`
	Service          ServiceConf    `toml:"service-conf" comment:"本服务配置"`
	TMClusterMonitor ClusterMonitor `toml:"tm-cluster-monitor" comment:"bsc集群监视"`
}

func LoadConf() *Config {
	data, _ := os.ReadFile("./config.toml")
	err := toml.Unmarshal(data, Conf)
	if err != nil {
		fmt.Println("get config info err,to exit server!,err is :%v", err)
		os.Exit(0)
	}
	return Conf
}

func (conf Config) String() string {
	b, _ := toml.Marshal(conf)
	return string(b)
}
