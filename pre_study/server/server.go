package main

import (
	"GoIn/pre_study/proto"
	"bufio"
	"fmt"
	"io"
	"net"
)

// TCP server端

// 处理函数
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	reader := bufio.NewReader(conn)
	for {
		// receive
		msg, err := proto.Decode(reader) // 读取数据
		if err == io.EOF {
			fmt.Println("客户端关闭链接!")
			break
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client端发来的数据：", msg)

		//send
		data, err := proto.Encode("server:"+msg)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			continue
		}
		conn.Write(data) // 发送数据
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}