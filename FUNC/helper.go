package netCat

import (
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	Clients  = make(map[net.Conn]Client)
	messages = []string{}
	mutex    = sync.Mutex{}
)
func valid(name string) bool {
	for _, i := range name {
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", i) {
			if (i == '\x1b' || i == '[') {
				return false
			}
		}
	}
	return true
}

func linuxLogo() string {
	return `
         nnnn
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
)      \.__.,|     .'
\__   )MMMMMP|   .'
	`
}
