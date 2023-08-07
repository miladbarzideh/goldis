package main

import (
	"bufio"
	"github.com/miladbarzideh/goldis/internal/network"
	"io"
	"log"
	"net"
)

func main() {
	socket, err := network.NewSocket(net.ParseIP("0.0.0.0"), 6380)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	for {
		c, err := socket.Accept()
		if err != nil {
			log.Panic(err)
		}

		//only serves one client connection at once
		for {
			err := oneRequest(c)
			if err != nil {
				log.Print(err)
				break
			}
		}
	}
}

func oneRequest(c *network.Socket) error {
	data, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("Client says: %s", data)

	if _, err = io.WriteString(c, "OK\n"); err != nil {
		return err
	}
	return nil
}
