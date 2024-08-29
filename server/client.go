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
	name = name[:len(name)-1]

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

	// Envoyer l'historique des messages aux nouveaux arrivant
	s.Mu.Lock()
	for _, msg := range s.History {
		var formattedMessage string
		var nameWithColor string
		if client.ID == msg.sender.ID {
			nameWithColor = Colorize("You", client.color)
			formattedMessage = fmt.Sprintf("[%s][%s]: %s", msg.timestamp.Format(time.RFC1123), nameWithColor, msg.content)
		} else {
			nameWithColor = Colorize(msg.sender.name, msg.sender.color)
			formattedMessage = fmt.Sprintf("[%s][%s]: %s", msg.timestamp.Format(time.RFC1123), nameWithColor, msg.content)
		}

		_, err := conn.Write([]byte(formattedMessage + "\n"))
		if err != nil {
			fmt.Println("Error sending message to new client:", err)
		}
	}
	s.Mu.Unlock()

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

		if content > 0 {
			message := Message{
				sender:    client,
				timestamp: time.Now(),
				content:   string(buffer[:content]),
			}

			if strings.HasPrefix(message.content, "/nick ") {
				newName := strings.TrimSpace(strings.TrimPrefix(message.content, "/nick "))
				if len(newName) > 0 {
					oldName := client.name
					client.name = newName
					s.MessageChannel <- Message{
						sender:    client,
						timestamp: time.Now(),
						content:   fmt.Sprintf("%s is now know has %s\n", oldName, newName),
					}
				}
				continue
			}

			if strings.HasPrefix(message.content, "/color ") {
				newColor := strings.TrimSpace(strings.TrimPrefix(message.content, "/color "))
				if len(newColor) > 0 {
					oldColor := client.color
					client.color = newColor

					oldColorName := Colorize(client.name, oldColor)
					newColorName := Colorize(client.name, newColor)

					s.MessageChannel <- Message{
						sender:    client,
						timestamp: time.Now(),
						content:   fmt.Sprintf("%s swith color %s to %s\n", client.name, oldColorName, newColorName),
					}
				}
				continue
			}

			if strings.HasPrefix(message.content, "/help") {
				helpMessage := `Available commands:
    /nick <new_name> - Change your nickname
    /color <color_name> - Change your text color
    /help - Show this help message`
				client.conn.Write([]byte(helpMessage + "\n"))
				continue
			}

			s.MessageChannel <- message
		}
	}
}

func GenerateUniqueID() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rng.Intn(10000))
}
