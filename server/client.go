package netcat

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strings"
	"time"
)

// Fonction qui gère la communication entre un client et le serveur
func (s *Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	if err := SendWelcomeMessage(conn); err != nil {
		fmt.Println(err)
		return
	}

	// Lecture du nom du client
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	name = GenerateFullName(name)
	color := randomColor()

	// Création du client
	client := Client{
		conn:  conn,
		name:  name,
		ID:    GenerateUniqueID(),
		color: color,
	}
	s.NewClientChannel <- client

	// Envoyer l'historique des messages aux nouveaux arrivants
	s.sendHistoryToClient(client)

	buffer := make([]byte, 1024)
	for {
		content, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				s.RemoveClientChannel <- conn
				fmt.Printf("%s has left the chat\n", client.name)
				return
			}
			fmt.Println("Error reading from client:", err)
			return
		}

		messageContent := strings.TrimSpace(string(buffer[:content]))
		if len(messageContent) == 0 {
			continue
		}

		// Gérer les commandes
		if s.handleCommand(messageContent, &client) {
			continue
		}

		// Envoyer le message
		message := Message{
			sender:    client,
			timestamp: time.Now(),
			content:   messageContent,
		}
		s.MessageChannel <- message
	}
}

// Fonction pour envoyer l'historique des messages à un nouveau client
func (s *Server) sendHistoryToClient(client Client) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, msg := range s.History {
		formattedMessage := formatMessageForClient(msg, client)
		_, err := client.conn.Write([]byte("\033[A\033[2K\r" + formattedMessage + "\n"))
		if err != nil {
			fmt.Println("Error sending message to new client:", err)
		}
	}
}

// Fonction pour gérer les commandes spécifiques
func (s *Server) handleCommand(messageContent string, client *Client) bool {
	if strings.HasPrefix(messageContent, "/nick ") {
		s.handleNickChange(strings.TrimSpace(strings.TrimPrefix(messageContent, "/nick ")), client)
		return true
	}

	if strings.HasPrefix(messageContent, "/color ") {
		s.handleColorChange(strings.TrimSpace(strings.TrimPrefix(messageContent, "/color ")), client)
		return true
	}

	if strings.HasPrefix(messageContent, "/help") {
		s.sendHelpMessage(client)
		return true
	}

	return false
}

// Fonction pour gérer le changement de pseudonyme
func (s *Server) handleNickChange(newName string, client *Client) {
	if len(newName) > 0 {
		parts := strings.SplitN(client.name, " ", 2)
		adjectives := ""

		if len(parts) == 2 {
			adjectives = parts[1]
		}

		newFullName := fmt.Sprintf("%s %s", newName, adjectives)

		oldName := client.name
		client.name = newFullName
		s.MessageChannel <- Message{
			sender:    *client,
			timestamp: time.Now(),
			content:   fmt.Sprintf("%s is now known as %s", oldName, newFullName),
		}
	}
}

// Fonction pour gérer le changement de couleur
func (s *Server) handleColorChange(newColor string, client *Client) {
	if len(newColor) > 0 {
		oldColor := client.color
		client.color = newColor

		oldColorName := Colorize(client.name, oldColor)
		newColorName := Colorize(client.name, newColor)

		s.MessageChannel <- Message{
			sender:    *client,
			timestamp: time.Now(),
			content:   fmt.Sprintf("%s switched color from %s to %s", client.name, oldColorName, newColorName),
		}
	}
}

// Fonction pour envoyer un message d'aide au client
func (s *Server) sendHelpMessage(client *Client) {
	helpMessage := `Available commands:
	/nick <new_name> - Change your nickname
	/color <color_name> - Change your text color
	/help - Show this help message`
	client.conn.Write([]byte("\033[A\033[2K\r" + helpMessage + "\n"))
}

// Fonction pour formater un message pour un client
func formatMessageForClient(msg Message, client Client) string {
	var nameWithColor string
	if client.ID == msg.sender.ID {
		nameWithColor = Colorize("You", client.color)
	} else {
		nameWithColor = Colorize(msg.sender.name, msg.sender.color)
	}
	return fmt.Sprintf("[%s][%s]: %s", msg.timestamp.Format(time.RFC1123), nameWithColor, msg.content)
}

func GenerateUniqueID() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rng.Intn(10000))
}
