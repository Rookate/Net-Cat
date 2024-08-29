package main

import (
	"flag"
	"fmt"
	"net"
	nc "netcat/server"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "Ip adress to bind to")
	port := flag.String("port", "8989", "Port to listen on")
	flag.Parse()

	logFile := nc.GenerateLogFileName()

	err := nc.OpenLogFile(logFile)

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

/*
	- Fonction pour gérer les connections et les déconnections des clients
	- Fonction pour gérer un seul client pour quand il envoie un message ou quand il recoit un message -> Goroutine pour chaque client
	- Fonction pour gérer les messages pour pouvoir les montrer à chaque clients. -> Utilisations des channels ainsi que mu pour synchroniser l'accès aux channels
*/
