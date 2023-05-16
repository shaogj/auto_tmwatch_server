
### 编译
go版本需>=1.18

`go build -o tmwatch_client tmwatch_client_triscript.go`

### 运行
`./tmwatch_client`

### 触发tm节点快照数据恢复接口
```bash
curl --location --request POST '127.0.0.1:6667/sync_tm_snapdata' \
--header 'Content-Type: application/json' \
--data-raw '{"auto_ip":"192,135","optype":"restoredata","snap_data_time":"20230513","token":"4444"}'

```

### 日志

日志每次触发，生成一个:前缀client_recoverrun_{日期时间}.log的新文件.示例为:`client_recoverrun20230513_{time}.log`

### 配置
配置文件为当前路径下的`config.toml`文件.北京环境的示例,有6个监控节点.为5个验证者+1个区块同步节点

