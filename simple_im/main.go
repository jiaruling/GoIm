package main

import s "GoIn/simple_im/server"

func main() {
	server := s.NewServer("0.0.0.0", 9091)
	server.Start()
}