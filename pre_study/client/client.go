package main

import (
	"GoIn/pre_study/proto"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// TCP Client端

// 客户端
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close() // 关闭连接
	reader := bufio.NewReader(conn)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// send
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, "\r\n") // 去除输入内容的前后空格
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			return
		}
		data, err := proto.Encode(inputInfo)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			continue
		}
		_, err = conn.Write(data) // 发送数据
		if err != nil {
			return
		}

		// receive
		msg, err := proto.Decode(reader)
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到server端发来的数据:",msg)
	}
}