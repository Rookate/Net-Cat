package netcat

import (
	"net"
	"sync"
	"time"
)

// Représente un client connecté au serveur
type Client struct {
	ID    string
	conn  net.Conn
	name  string
	color string
}

// Server gère les clients et les messages
type Server struct {
	Clients             map[net.Conn]Client
	Mu                  sync.Mutex
	MessageChannel      chan Message
	NewClientChannel    chan Client
	RemoveClientChannel chan net.Conn
	History             []Message
}

type Message struct {
	timestamp time.Time
	sender    Client
	content   string
}

// NewServer crée une nouvelle instance de Server
func NewServer() *Server {
	return &Server{
		Clients:             make(map[net.Conn]Client),
		MessageChannel:      make(chan Message),
		NewClientChannel:    make(chan Client),
		RemoveClientChannel: make(chan net.Conn),
	}
}
