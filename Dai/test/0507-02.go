package main

import "net"

func main() {
	server, err := net.Listen("tcp", "127.0.0.1")
	if err != nil {
		panic(err)
	}

}
