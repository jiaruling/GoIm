# Go语言实战: 即使通信系统

## 参考

**[简易通讯系统](https://www.cnblogs.com/N3ptune/p/16268670.html)**

## 需求

- 公聊
- 私聊
- 查看在线用户
- 用户强制下线

## 运行

```shell script
# 启动服务器
~: go run main.go
# 启动客户端
~: cd client && go run client.go
# 通过linux命令连接服务器
~: nc 192.168.8.88 9091
```