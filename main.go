package main

import (
	"flag"
	"fmt"
	"net"
	nc "netcat/server"
)

func main() {
	ip := flag.String("ip", "172.20.10.3", "Ip adress to bind to")
	port := flag.String("port", "8989", "Port to listen on")
	flag.Parse()

	logFile := nc.GenerateLogFileName()

	err := nc.OpenLogFile(logFile)

	if err != nil {
		fmt.Println(err)
	}

	defer nc.CloseLogFile()

	adress := fmt.Sprintf("%s:%s", *ip, *port)

	server := nc.NewServer()

	listener, err := net.Listen("tcp", adress)

	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Server is listening on:", adress)
	go server.HandleMessage()
	go server.ManageClient()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accecpting connection", err)
			continue
		}
		go server.HandleClient(conn)
	}
}
