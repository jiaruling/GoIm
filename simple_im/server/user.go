package server

import (
	"context"
	"net"
	"strings"
)

type User struct {
	Name string      // 名称
	Addr string      // 地址
	ch   chan string // 通道 用来向客户端发送消息
	conn net.Conn    // 连接
	serv *Server
}

// NewUser 创建一个用户
func NewUser(conn net.Conn, serv *Server, ctx context.Context) *User {
	userAddr := conn.RemoteAddr().String()
	// 创建结构体
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		ch:   make(chan string),
		conn: conn,
		serv: serv,
	}
	go user.ListenMessage(ctx) // 调用一个goroutine监听消息
	return user
}

// ListenMessage 监听当前user的通道
func (user *User) ListenMessage(ctx context.Context) {
loop:
	for {
		select {
		case msg := <-user.ch:
			user.conn.Write([]byte(msg + "\n")) // 发送给客户端
		case <-ctx.Done():
			break loop
		}
	}
	return
}

// Online 用户上线
func (user *User) Online() {
	user.serv.mapLock.Lock()
	user.serv.OnlineMap[user.Name] = user
	user.serv.mapLock.Unlock()
	user.serv.BroadCast(user, " has arrived")
}

// Offline 用户下线
func (user *User) Offline() {
	// 从map去除
	user.serv.mapLock.Lock()
	delete(user.serv.OnlineMap, user.Name)
	user.serv.mapLock.Unlock()
	// 广播上线消息
	user.serv.BroadCast(user, " has left")
}

func (user *User) DoMessage(msg string, closeClient chan struct{}, exit *bool) {
	// 执行指令
	if msg[0] == '$' {
		if msg == "$who" {
			// 查询在线用户
			user.ch <- "Online users:"
			user.serv.mapLock.Lock()
			for _, u := range user.serv.OnlineMap {
				user.ch <- "[" + u.Addr + "] " + u.Name
			}
			user.serv.mapLock.Unlock()
		} else if len(msg) > 7 && msg[:7] == "$rename" {
			// 修改用户名称
			name := strings.Split(msg, ":")[1] // 字符串切割，以:为分界
			_, ok := user.serv.OnlineMap[name] // 检查name是否已经存在
			if ok {
				user.ch <- "User name already exists"
			} else {
				user.serv.mapLock.Lock()
				delete(user.serv.OnlineMap, user.Name) // 删除原先的name
				user.serv.OnlineMap[name] = user       // 添加新的name
				user.ch <- "Your name has been updated: " + name
				user.Name = name
				user.serv.mapLock.Unlock()
			}
		} else if len(msg) > 4 && msg[:4] == "$to " {
			name := strings.Split(msg, ":")[0][4:] // 接收者
			send := strings.Split(msg, ":")[1]     // 消息内容
			receiver, ok := user.serv.OnlineMap[name]
			if ok {
				receiver.ch <- user.Name + " to you: " + send
				user.ch <- "send ok"
			} else {
				user.ch <- "The user does not exist"
			}
		} else if msg == "$exit" {
			closeClient <- struct{}{}
			*exit = true
		} else {
			// 提示信息
			user.ch <- "Wrong instruction"
			user.ch <- "$who:view online users\n$rename:modify name"
		}
		return
	}
	user.serv.BroadCast(user, msg)
}
