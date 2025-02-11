package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients  = make(map[net.Conn]Client)
	messages = []string{}
	mutex    = sync.Mutex{}
)

func main() {
	port := ":8989"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer ln.Close()
	fmt.Println("Listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		if len(clients) >= 10 {
			conn.Write([]byte("Server full. Try again later.\n"))
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conten, _ := os.ReadFile("fil.txt")

	conn.Write([]byte("Welcome to TCP-Chat!\n" + string(conten) + "\n[ENTER YOUR NAME]: "))
	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		conn.Write([]byte("Invalid name. Connection closed.\n"))
		return
	}

	mutex.Lock()
	clients[conn] = Client{conn, name}
	mutex.Unlock()

	broadcast(fmt.Sprintf("\n%s has joined the chat", name), conn)
	conn.Write([]byte(strings.Join(messages, "\n") + "\n"))

	for {
		formattedMsg0 := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)

		conn.Write([]byte(formattedMsg0 + " "))

		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		formattedMsg := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, msg)
		mutex.Lock()
		messages = append(messages, formattedMsg)
		mutex.Unlock()
		broadcast(formattedMsg, conn)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	broadcast(fmt.Sprintf("%s has left the chat...", name), conn)
}

func broadcast(msg string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(msg + "\n"))
		}
	}
}
