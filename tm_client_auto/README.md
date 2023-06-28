
# tmnode_watching_status

此项目用于重启tmware的服务节点,需赋予sudo的权限

项目支持通过当前目录下的config.toml来配置

## 配置
配置文件为当前路径下的`config.toml`文件.北京环境的示例,有6个监控节点.为5个验证者+1个区块同步节点
### 配置文件config.toml参数说明
```
[service-conf]

  //服务端口默认为6667
  port = 6667
  //日志等级int类型默认为info; debug || info || warn || error
  log-level = "info"
  //日志路径,默认为update.log
  log-path = "./tmwatch_client.log"

 //用于触发tmclient请求的access token
  access-token = "3333"
  //调用执行的tm脚本文件名
  invoke-script-name = "tmnode_snapdata0419KK.sh"


//检测节点状态配置
[service-conf-monitor]
  //检测节点接口间隔时间
  request_tm_interval = 60
  //重启节点服务间隔时间
  operator_sys_interval = 30
  //节点服务启动方式
  start_mode = 0
```

### 编译
go版本需>=1.18

`go build -o tmwatch_client tmwatch_client_triscript.go
或编译linux平台应用，使用指令：
make linux-amd64
`

### 运行
`./tmwatch_client
或Linux平台：
./tmnode_watching_status
`

### 触发tm节点快照数据恢复接口
```bash
curl --location --request POST '127.0.0.1:6667/sync_tm_snapdata' \
--header 'Content-Type: application/json' \
--data-raw '{"auto_ip":"192,135","optype":"restoredata","snap_data_time":"20230513","token":"4444"}'

```

### 日志

日志每次触发，生成一个:前缀client_recoverrun_{日期时间}.log的新文件.示例为:`client_recoverrun20230513_{time}.log`

