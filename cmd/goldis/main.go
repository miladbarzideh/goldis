package main

import (
	"log"
	"net"

	"github.com/miladbarzideh/goldis/internal/network"
)

const (
	ip   = "0.0.0.0"
	port = 6380
)

func main() {
	socket, err := network.NewSocket(net.ParseIP(ip), port)
	if err != nil {
		log.Fatal(err)
	}

	defer socket.Close()

	connManager := network.NewConnectionHandler(socket)
	connManager.StartServer()
}
