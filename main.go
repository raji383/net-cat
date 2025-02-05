package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	clients  = make(map[net.Conn]string)
	messages = make(chan string)
	mutex    = sync.Mutex{}
)

func main() {
	port := "8989"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Listening on port:", port)

	go broadcastMessages()

	for {
		if len(clients) >= 10 {
			fmt.Println("Maximum clients reached.")
			continue
		}

		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	contint,_ := os.ReadFile("fil.txt")
	fmt.Fprint(conn, "Welcome to TCP-Chat!\n")
	fmt.Fprint(conn,string(contint))
	fmt.Fprint(conn, "[Enter your name]: ")

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	name := scanner.Text()
	if strings.TrimSpace(name) == "" {
		fmt.Fprintln(conn, "Name cannot be empty. Closing connection.")
		return
	}

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	messages <- fmt.Sprintf("%s has joined the chat.", name)
	fmt.Printf("%s joined the chat.\n", name)

	for scanner.Scan() {
		msg := scanner.Text()
		if strings.TrimSpace(msg) == "" {
			continue
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		messages <- fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	messages <- fmt.Sprintf("%s has left the chat.", name)
	fmt.Printf("%s left the chat.\n", name)
}

func broadcastMessages() {
	for msg := range messages {
		mutex.Lock()
		for conn := range clients {
			fmt.Fprintln(conn, msg)
		}
		mutex.Unlock()
	}
}
