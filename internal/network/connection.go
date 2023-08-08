package network

import (
	"fmt"
	"syscall"
)

var (
	MAXSIZE = 8000
)

type Connection struct {
	Fd   int
	Addr syscall.Sockaddr
}

func (c Connection) Read() (string, error) {
	msg := make([]byte, MAXSIZE)
	sizeMsg, _, err := syscall.Recvfrom(c.Fd, msg, 0)
	if err != nil {
		return "", err
	}
	addrFrom := c.Addr.(*syscall.SockaddrInet4)
	message := string(msg)

	fmt.Printf("%d byte read from %d:%d on socket %d\n", sizeMsg, addrFrom.Addr, addrFrom.Port, c.Fd)
	fmt.Printf("Received message: %s\n", message)

	return message, nil
}

func (c Connection) Write(msg []byte) (int, error) {
	err := syscall.Sendmsg(c.Fd, msg, nil, c.Addr, syscall.MSG_DONTWAIT)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Response message: %s ", msg)
	return 1, err
}

func (c Connection) Close() error {
	return syscall.Close(c.Fd)
}
