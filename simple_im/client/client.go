package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	ServIP   string   // 服务器IP
	ServPort int      // 端口号
	Name     string   // 名称
	conn     net.Conn // 连接
	choice   int
}

var (
	addr string
	port int
)

func (clnt *Client) menu() bool {
	var choice int

	fmt.Println("******************************")
	fmt.Println("1.Public chat mode")
	fmt.Println("2.Private chat mode")
	fmt.Println("3.Update the user name")
	fmt.Println("0.Exit")
	fmt.Println("******************************")

	fmt.Scanf("%d", &choice)
	if choice >= 0 && choice <= 3 {
		clnt.choice = choice
		return true
	} else {
		fmt.Println("Illegal input")
		return false
	}
}

// UpdateName 更新用户名
func (clnt *Client) UpdateName() bool {
	fmt.Println("Enter your new username: ")
	// 接收控制台输入
	clnt.Name = input()
	if clnt.Name == "" {
		fmt.Println("获取名称失败")
		return false
	}
	msg := "$rename:" + clnt.Name + "\n"
	_, err := clnt.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// PublicChat 公聊
func (clnt *Client) PublicChat() {
	// 提示用户输入消息
	var msg string
	fmt.Println("[Public Mode]\nEnter '$exit' to exit")
	reader := bufio.NewReader(os.Stdin)
	for msg != "$exit" {
		if len(msg) > 0 {
			send := msg + "\n"
			_, err := clnt.conn.Write([]byte(send))
			if err != nil {
				log.Println(err)
				break
			}
		}
		msg = ""
		text, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(text)
		//fmt.Scanln(&msg)
	}
}

// 客户端关闭通知服务器
func (clnt *Client) Close() {
	send := "$exit\n"
	_, err := clnt.conn.Write([]byte(send))
	if err != nil {
		log.Println(err)
		return
	}
}

// 显示所有在线用户
func (clnt *Client) SelectUser() {
	send := "$who\n"
	_, err := clnt.conn.Write([]byte(send))
	if err != nil {
		log.Println(err)
		return
	}
}

// PrivateChat 私聊
func (clnt *Client) PrivateChat() {
	var msg string
	fmt.Println("[Private Mode] Enter '$exit' to exit")
	clnt.SelectUser()
	// 接收控制台输入
	name := input()
	if name == "" {
		fmt.Println("获取名称失败")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for msg != "$exit" {
		if len(msg) > 0 {
			send := "$to " + name + ":" + msg + "\n"
			_, err := clnt.conn.Write([]byte(send))
			if err != nil {
				log.Println(err)
				break
			}
		}
		msg = ""
		text, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(text)
		//fmt.Scanln(&msg)
	}
}

func input() (value string) {
	for {
		_, err := fmt.Scanln(&value)
		if err == nil {
			break
		} else if err.Error() == "unexpected newline" {
			//fmt.Println("continue")
			continue
		} else {
			fmt.Println(err.Error())
			return
		}
	}
	return
}

func (clnt *Client) Run() {
	for clnt.choice != 0 {
		time.Sleep(200 * time.Millisecond)
		for !clnt.menu() {
		} // 如果不为true，则一直循环在这里
		// 根据不同的模式处理不同业务
		switch clnt.choice {
		case 1: // 公聊模
			clnt.PublicChat()
		case 2: // 私聊模式
			clnt.PrivateChat()
		case 3: // 更新用户名
			clnt.UpdateName()
		case 0: // 为0则循环结束
			fmt.Println("exit...")
			clnt.Close()
		}
	}
}

// DoResponse 处理服务器消息
func (clnt *Client) DoResponse() {
	io.Copy(os.Stdout, clnt.conn)
}

func init() {
	// flag.TypeVar(Type 指针, flag 名, 默认值, 帮助信息)
	flag.StringVar(&addr, "IP", "127.0.0.1",
		"Set the server IP address")
	flag.IntVar(&port, "port", 9091,
		"Set the server port number")
}

func NewClient(addr string, port int) *Client {
	// 创建客户端对象
	clnt := &Client{
		ServIP:   addr,
		ServPort: port,
		choice:   999,
	}
	// 发起请求 拼接地址和端口
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Println(err)
		return nil
	}
	clnt.conn = conn
	return clnt
}

func main() {
	flag.Parse()
	clnt := NewClient(addr, port)
	if clnt == nil {
		return
	}
	go clnt.DoResponse()
	fmt.Println("connect ok")
	clnt.Run()
	//select {}
}
