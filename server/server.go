package netcat

import (
	"fmt"
	"time"
)

// Fonction qui va envoyer un message à chaque client sauf à l'emetteur
func (s *Server) HandleMessage() {
	for message := range s.MessageChannel {
		s.Mu.Lock()
		s.History = append(s.History, message)
		s.Mu.Unlock()

		logMessage := fmt.Sprintf("[%s][%s]: %s", message.timestamp.Format(time.RFC1123), message.sender.name, message.content)

		LogMessage(logMessage)

		for _, client := range s.Clients {
			var formattedMessage string
			var nameWithColor string
			if client.ID == message.sender.ID {
				// Message envoyé par le client lui-même
				nameWithColor = Colorize("You", message.sender.color)
				formattedMessage = fmt.Sprintf("[%s][%s]: %s", message.timestamp.Format(time.RFC1123), nameWithColor, message.content)
			} else {
				// Message envoyé par un autre client
				nameWithColor = Colorize(message.sender.name, message.sender.color)
				formattedMessage = fmt.Sprintf("[%s][%s]: %s", message.timestamp.Format(time.RFC1123), nameWithColor, message.content)
			}
			_, err := client.conn.Write([]byte(formattedMessage + "\n"))
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}

func (s *Server) ManageClient() {
	for {
		select {
		case newClient := <-s.NewClientChannel:
			s.Mu.Lock()
			s.Clients[newClient.conn] = newClient
			s.Mu.Unlock()

			notification := fmt.Sprintf("%s has joined the chat at %s\n", newClient.name, time.Now().Format(time.TimeOnly))

			LogMessage(notification)
			fmt.Print(notification)
			for _, client := range s.Clients {
				client.conn.Write([]byte(notification))
			}
		case disconnectedClient := <-s.RemoveClientChannel:
			s.Mu.Lock()
			client := s.Clients[disconnectedClient]
			delete(s.Clients, disconnectedClient)
			s.Mu.Unlock()

			notification := fmt.Sprintf("%s has left the chat at %s\n", client.name, time.Now().Format(time.TimeOnly))
			LogMessage(notification)

			for _, client := range s.Clients {
				client.conn.Write([]byte(notification))
			}
		}
	}
}
