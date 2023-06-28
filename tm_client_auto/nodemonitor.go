package main

import (
	//"2021New_BFLProjTotal/tm_watch_projmy/config"
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/config"
	"202108FromBFLProj/auto_tmwatch_server/tm_client_auto/unit_req"
	"errors"
	"fmt"
	"github.com/mkideal/log"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return ip, err

}

func GetExternIP2() {
	responseClient, errClient := http.Get("http://ip.dhcp.cn/?ip") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		panic(errClient)
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()

	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	print(clientIP)
}

// 0110add reset proc
func IsProcessExist(appName string) bool {

	cmd := exec.Command("ps", "-C", appName)
	output, _ := cmd.Output()

	fields := strings.Fields(string(output))

	for _, v := range fields {
		if v == appName {
			return true
		}
	}

	return false
}

// start process
func ExecProcess(appName string) {

	path := "./" // app路径

	cmd := exec.Command(path + appName)
	cmd.Output()
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func StartMonitorTM(tmMonitorCfg config.ConfigMonitorInfomationTM) {
	getlocalIp, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("getlocal2023 ip is :%s", getlocalIp)

	if err := config.InitWithProviders("multifile/console", "./logs"); err != nil {
		panic("init log error: " + err.Error())
	}
	log.Info("get cur ip: %s,run StartTMWatchLocalServer() !", getlocalIp)

	log.Info("log level: %v", log.SetLevelFromString("trace"))
	/*
		err = config.InitValidatorConfigInfo()
		if nil != err {
			fmt.Println("test InitValidatorConfigInfo()err! ,err is:%v", err)
			log.Error("from config.json,get json conf err!,err is:%v", err)
			os.Exit(0)
		}
	*/
	//tmMonitorCfg := &config.GbTrustConf
	//tmMonitorCfg := &tmMonitorCfg
	log.Info("2023----doing----monitor tmnode info: %v", tmMonitorCfg)

	ch := make(chan int)
	go func() {
		//0110,check,tmnode status
		//unit_req.MonitorTMBlockIncrease(gbConf, getlocalIp.String())
		unit_req.MonitorTMBlockIncrease(&tmMonitorCfg, getlocalIp.String())
	}()

	time.Sleep(time.Duration(3) * time.Second)

	curpath := unit_req.GetAppPath()
	log.Info("checking=== cur path is:%s\n", curpath)

	<-ch
	//select {}
	fmt.Println("checking now to MonitorBlockIncrease()!")

}
