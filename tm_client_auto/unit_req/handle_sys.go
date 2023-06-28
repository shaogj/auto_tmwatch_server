package unit_req

import (
	"fmt"
	"github.com/mkideal/log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func GetAppPath() string {
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	dir := filepath.Dir(path)

	return dir
}

// use:	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
// cmd := exec.Command("/bin/bash", "-c", "cd /app/read5 && ./read5 serve -d")
// runInLinux（）
func KillFindPid(strpid string) error {
	ipid, _ := strconv.Atoi(strpid)
	var err error
	if strpid != "" {
		log.Info("in KillFindPid(),cur to kill pid is:%s\n", strpid)
		err = syscall.Kill(ipid, syscall.SIGKILL)
	} else {
		log.Error("in KillFindPid(),no exec no Kill proc! exist no the pid is :%s \n", strpid)
	}
	return err
}

func KillProcessByName(procname string) error {
	var pid string
	var err error
	if runtime.GOOS == "darwin" {
		pid, err = GetPid(procname)
	} else {
		pid, err = GetPid(procname)
	}
	log.Info("checking===In KillProcessByName(),cur to kill procname is:%s,pid is:%v,err is:%v\n", procname, pid, err)
	time.Sleep(time.Duration(2) * time.Second)
	error := KillFindPid(pid)
	return error
}
func runInLinux(cmd string) (string, error) {
	log.Info("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return "", nil //(cmd)
	} else {
		return runInLinux(cmd)
	}
}

// 根据进程名称获取进程ID
func GetPid(serverName string) (string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	//add

	return pid, err
}

//0112add:sudo systemctl start tendermint
func StartTendermint(serverName string) (string, error) {
	//a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	acmd := `sudo systemctl start tendermint`
	pid, err := RunCommand(acmd)

	return pid, err
}
func StopTendermint(serverName string) (string, error) {
	acmd := `sudo systemctl stop tendermint`
	pid, err := RunCommand(acmd)

	return pid, err
}

func StartTMServerLocal(serverName string) (string, error) {
	//good!
	cmdstarttm := `nohup /Users/gejians/go/src/tmware/cmd/tendermint/tendermint --home /Users/gejians/go/src/tmware/cmd/tendermint/test2023new1 start --log-level debug --proxy-app=persistent_kvstore &`
	pid, err := RunCommand(cmdstarttm)

	return pid, err
}

func StartTMServer() (string, error) {
	cmdstarttm := `nohup /usr/local/bin/tendermint --home /trias/.ethermint/tendermint start --log-level debug --proxy-app=persistent_kvstore --consensus.create-empty-blocks=false &`
	pid, err := RunCommand(cmdstarttm)

	return pid, err
}
