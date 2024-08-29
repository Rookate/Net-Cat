package netcat

import (
	"fmt"
	"math/rand"
	"net"
)

func GenerateFullName(name string) string {
	names := []string{"Alice", "Bob", "Charles", "Pcaboo", "Adridri", "Jhon"}
	adjectives := []string{"Le Malicieux", "Le Gourgandin", "Le Terrible", "Le Dépressif"}

	adjIndex := rand.Intn(len(adjectives))

	var nameIndex int
	var finalName string

	tag := rand.Intn(10000)

	if len(name) == 0 {
		nameIndex = rand.Intn(len(names))
		finalName = fmt.Sprintf("%s %s#%4d", names[nameIndex], adjectives[adjIndex], tag)
	} else {
		finalName = fmt.Sprintf("%s %s", name, adjectives[adjIndex])
	}

	return finalName
}

func SendWelcomeMessage(conn net.Conn) error {
	welcomeMessage := `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '\ \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'

Enter your name:
`
	_, err := conn.Write([]byte(welcomeMessage))
	if err != nil {
		return fmt.Errorf("error sending welcome message: %v", err)
	}
	return nil
}
