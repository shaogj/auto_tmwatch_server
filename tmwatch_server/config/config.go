package config

import (
	log "202108FromBFLProj/ChainWatch_Project2023/auto_proj_test/logs"
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"sync"
	"time"
)

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
	DingUrl       string   `toml:"ding-url" comment:"钉钉机器人接口"`
	OkPrefixKey   string   `toml:"ok-prefix-key" comment:"钉钉安全设置,节点正常消息前缀"`
	ErrPrefixKey  string   `toml:"err-prefix-key" comment:"钉钉安全设置,节点异常消息前缀"`
	Interval      int      `toml:"interval" comment:"获取块高时间间隔,单位为秒"`
	RetryTimes    int      `toml:"retry-times" comment:"判断区块落后所需次数"`
	AbnormalHosts []string `toml:"abnormal-hosts" comment:"落后区块节点ip列表"`
}

type ServiceConf struct {
	Port     int    `toml:"port" comment:"服务端口默认为6667"`
	LogLevel string `toml:"log-level" comment:"日志等级int类型默认为info; debug || info || warn || error"`
	LogPath  string `toml:"log-path" comment:"日志路径,默认为update.log"`
}

type ClusterMonitor struct {
	ClusterInterval int `toml:"cluster-interval" comment:"查看集群状态时间间隔,默认10分钟"`
	NodeInterval    int `toml:"node-interval" comment:"bsc节点重启时间间隔,默认2分钟"`
	MonitorRpc      int `toml:"monitor-rpc" comment:"用于重启bsc的rpc端口"`
}
type Config struct {
	TM                []Host         `toml:"tm" comment:"tm链节点ip列表"`
	BSC               []Host         `toml:"bsc" comment:"bsc链节点ip列表"`
	BscServer         BscConf        `toml:"bsc-conf" comment:"bsc链server配置"`
	TmServer          TmConf         `toml:"tm-conf" comment:"tm链server ip"`
	TmMonitor         Monitor        `toml:"tm-monitor" comment:"tm落后节点监视器"`
	BscMonitor        Monitor        `toml:"bsc-monitor" comment:"bsc节点停止出块监视"`
	Service           ServiceConf    `toml:"service-conf" comment:"本服务配置"`
	BscClusterMonitor ClusterMonitor `toml:"bsc-cluster-monitor" comment:"bsc集群监视"`
}

func init() {
	LoadConf()
}

func LoadConf() *Config {
	data, _ := os.ReadFile("./config.toml")
	err := toml.Unmarshal(data, Conf)
	fmt.Println(err)

	return Conf
}

func SaveConf(conf Config) error {
	log.Logger.Info("save conf...")
	defer log.Logger.Info("log saved")
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// 先备份conf文件
	var ConfName string = "./config.toml"
	BakConf := fmt.Sprintf("%sold-%s", time.Now().Format("2006-01-02_15_04"), ConfName[2:])
	err := os.Rename(ConfName, BakConf)
	if err != nil {
		return err
	}
	// open file以追加的方式写入,toml write 会存在乱码的问题
	// f, err := os.OpenFile("./config.toml", os.O_CREATE|os.O_WRONLY, 0644)
	// 可能需要像bsc一样复制
	f, err := os.Create("./config.toml")
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	enc.Order(toml.OrderPreserve)

	if err = enc.Encode(conf); err != nil {
		return err
	}
	return nil
}

// 获取去重后的ip地址列表
func GetIps(nodeType string) ([]string, error) {
	temp := make(map[string]bool)
	var ips []string
	if nodeType == "bsc" {
		for _, host := range Conf.BSC {
			if _, ok := temp[host.Ip]; !ok {
				temp[host.Ip] = host.IsNew
				// 不再支持从配置文件添加
				// if host.IsNew {
				// 	ips = append(ips, host.Ip)
				// }
				ips = append(ips, host.Ip)
			}
		}
	} else if nodeType == "tm" {
		for _, host := range Conf.TM {
			if _, ok := temp[host.Ip]; !ok {
				temp[host.Ip] = host.IsNew
				ips = append(ips, host.Ip)
			}
		}
	} else {
		err := fmt.Errorf("bad nodeType: %s", nodeType)
		log.Logger.Error(err)
		return nil, err
	}

	return ips, nil
}

func (conf Config) String() string {
	b, _ := toml.Marshal(conf)
	return string(b)
}
