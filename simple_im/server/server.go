package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP        string           // IP地址
	Port      int              // 端口号
	OnlineMap map[string]*User // 在线用户列表
	mapLock   sync.RWMutex     // 读写锁
	Message   chan string      // 广播消息
}

// NewServer 创建一个服务器
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// handler 处理连接
func (s *Server) handler(conn net.Conn) {
	var exit bool
	ctx, cancel := context.WithCancel(context.Background())
	addr := conn.RemoteAddr().String()
	log.Println(addr + " is connecting")
	// 创建用户
	user := NewUser(conn, s, ctx)
	// 用户上线
	user.Online()
	// 是否活跃
	active := make(chan struct{})
	// 客户端关闭
	closeClient := make(chan struct{})
	// 接收消息
	go func(ctx context.Context) {
		input := bufio.NewScanner(conn)
	loop:
		for {
			select {
			case <-ctx.Done():
				// 用户下线
				user.Offline()
				log.Println(addr + " has been disconnected")
				break loop
			default:
				if input.Scan() {
					log.Println(addr + ": " + input.Text())
					//s.BroadCast(user, input.Text())
					user.DoMessage(input.Text(), closeClient, &exit)
					if !exit {
						// 用户的任意消息 都代表用户当前活跃
						active <- struct{}{}
					}
				}
			}
		}
	}(ctx)

	for {
		select {
		case <-active:
			// 激活select便会 重置定时器
		case <-time.After(time.Minute * 10):
			// 超时
			s.exit(user, active, closeClient, "You've been kicked out\n")
			cancel()
			return
		case <-closeClient:
			// 客户端主动关闭
			s.exit(user, active, closeClient)
			cancel()
			return
		}
	}
}

func (s *Server) exit(user *User, active, closeClient chan struct{}, send ...string) {
	// 超时
	if len(send) > 0 {
		user.conn.Write([]byte(send[0]))
	}
	// 关闭通道
	close(user.ch)
	close(active)
	close(closeClient)
	user.conn.Close() // 关闭连接
	return            // 返回这个函数
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg // 向服务器广播消息通道发送字符串
}

// 广播消息
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		// 将msg发送给全部在线的User
		s.mapLock.Lock()
		for _, userClient := range s.OnlineMap {
			userClient.ch <- msg
		}
		s.mapLock.Unlock()
	}
}

// Start 启动服务器
func (s *Server) Start() {
	// 监听TCP端口
	listener, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close() // 关闭连接
	// 启动监听Message的goroutine
	go s.ListenMessage() // 广播消息
	for {
		// 接收连接
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handler(conn)
	}
}
