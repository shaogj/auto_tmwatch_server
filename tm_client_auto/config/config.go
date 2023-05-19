package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
)

// var Conf = &Config{}
var Conf = &Config{}

type Host struct {
	Ip    string `toml:"ip" comment:"ip地址"`
	IsNew bool   `toml:"isNew" comment:"用于确定是否为新加ip,默认为false,已弃用"`
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
	//0506add
	DingToken string `toml:"ding-token" comment:"钉钉机器人Token接口"`
	//ding-token = '48679669f51e7aa89c4544b17bfc4c1e4a565b975d067cf3fa420b0ae0c255ef'

}
type ServiceConf struct {
	Port             int    `toml:"port" comment:"服务端口默认为6667"`
	LogLevel         string `toml:"log-level" comment:"日志等级int类型默认为info; debug || info || warn || error"`
	LogPath          string `toml:"log-path" comment:"日志路径,默认为update.log"`
	AccessToken      string `toml:"access-token" comment:"访问tmwatch-client的请求token"`
	InvokeScriptName string `toml:"invoke-script-name" comment:"调用执行的tm脚本文件名"`
}

type Config struct {
	//TmServer          TmConf         `toml:"tm-conf" comment:"tm链server ip"`
	TmMonitor Monitor     `toml:"tm-monitor" comment:"tm落后节点监视器"`
	Service   ServiceConf `toml:"service-conf" comment:"本服务配置"`
}

func init() {
	//LoadConf()
}

func LoadConf() *Config {
	data, _ := os.ReadFile("./config.toml")
	err := toml.Unmarshal(data, Conf)
	if err != nil {
		fmt.Println("get config info err,to exit server!,err is :%v", err)
		os.Exit(0)
	}
	//fmt.Println("get config info is===77:%v:", Conf.Service)
	//log.Logger.Infof("get config.toml info' Conf.Service ========= is :%v", Conf.Service)
	//Logger
	return Conf
}

func (conf Config) String() string {
	b, _ := toml.Marshal(conf)
	return string(b)
}
