package netCat

import (
	"fmt"
	"net"
	"time"
)

func broadcast(msg string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn, name := range Clients {
		if conn != sender {
			formattedMsg := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name.name)
			conn.Write([]byte("\n" + msg + formattedMsg))
		}
	}
}
