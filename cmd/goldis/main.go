package main

import (
	"log"
	"net"

	"github.com/miladbarzideh/goldis/internal/network"
)

func main() {
	socket, err := network.NewSocket(net.ParseIP("0.0.0.0"), 6380)
	if err != nil {
		log.Fatal(err)
	}

	defer socket.Close()

	connManager := network.NewConnectionHandler(socket)
	connManager.StartServer()
}
