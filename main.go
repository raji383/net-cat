package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"netCat/FUNC"
)

func main() {
	port := ":8989"
	if len(os.Args) > 2 {
		return
	}
	if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	}
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
		if len(netCat.Clients) > 2 {
			conn.Write([]byte("Server full. Try again later.\n"))
			conn.Close()
			continue
		}
		go netCat.HandleConnection(conn)
	}
}
