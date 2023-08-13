package network

import (
	"log"
	"syscall"
)

const (
	maxSize = 1024
)

type Connection struct {
	Fd   int
	Addr syscall.Sockaddr
}

func (c Connection) Read() (string, error) {
	buf := make([]byte, maxSize)
	sizeMsg, _, err := syscall.Recvfrom(c.Fd, buf, 0)
	if err != nil {
		return "", err
	}

	input := string(buf[:sizeMsg])

	addrFrom := c.Addr.(*syscall.SockaddrInet4)
	log.Printf("%d byte read from %d:%d on socket %d\n", sizeMsg, addrFrom.Addr, addrFrom.Port, c.Fd)
	log.Printf("Received command: %s\n", input)

	return input, nil
}

func (c Connection) Write(msg []byte) (int, error) {
	err := syscall.Sendmsg(c.Fd, msg, nil, c.Addr, syscall.MSG_DONTWAIT)
	if err != nil {
		return 0, err
	}
	log.Printf("Response message: %s ", msg)
	return 1, err
}

func (c Connection) Close() error {
	return syscall.Close(c.Fd)
}
