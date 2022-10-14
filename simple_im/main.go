package main

import (
	s "GoIn/simple_im/server"
	"fmt"
	"runtime"
	"time"
)

func main() {
	// 测试, 关闭一个连接后是否关闭了所对应的goroutine
	go func() {
		for {
			fmt.Println("当前goroutine数量:",runtime.NumGoroutine())
			time.Sleep(3 * time.Second)
		}
	}()
	server := s.NewServer("0.0.0.0", 9091)
	server.Start()
}
