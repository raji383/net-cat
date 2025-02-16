package netCat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	conn.Write([]byte("Welcome to TCP-Chat!\n" + linuxLogo() + "\n[ENTER YOUR NAME]: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	for !valdName(name) {
		conn.Write([]byte("Invalid name. Connection closed." + "\n[ENTER YOUR NAME]: "))
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)

	}

	for {
		mutex.Lock()
		duplicate := false
		for _, client := range Clients {
			if client.name == name {
				duplicate = true
				break
			}
		}
		mutex.Unlock()

		if duplicate || name == "" {
			conn.Write([]byte("Name already taken" + "\n[ENTER YOUR NAME]: "))
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)

			continue
		}
		break
	}

	mutex.Lock()
	Clients[conn] = Client{conn, name}
	mutex.Unlock()

	broadcast(fmt.Sprintf("%s has joined the chat", name), conn)
	conn.Write([]byte(strings.Join(messages, "\n") + "\n"))

	for {
		formattedMsg0 := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)

		mutex.Lock()
		conn.Write([]byte(formattedMsg0 + " "))
		mutex.Unlock()

		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		mutex.Lock()
		messages = append(messages, formattedMsg0+msg)
		mutex.Unlock()
		broadcast(formattedMsg0+msg, conn)
	}

	mutex.Lock()
	delete(Clients, conn)
	mutex.Unlock()
	broadcast(fmt.Sprintf("%s has left the chat...", name), conn)
}
