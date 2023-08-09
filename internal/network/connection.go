package network

import (
	"log"
	"regexp"
	"syscall"
)

var (
	MAXSIZE = 1024
)

type Connection struct {
	Fd   int
	Addr syscall.Sockaddr
}

func (c Connection) Read() ([]string, error) {
	buf := make([]byte, MAXSIZE)
	sizeMsg, _, err := syscall.Recvfrom(c.Fd, buf, 0)
	if err != nil {
		return nil, err
	}

	// Ignore white spaces and new lines
	input := string(buf[:sizeMsg])
	regex := regexp.MustCompile(`\w+`)
	commands := regex.FindAllString(input, -1)

	addrFrom := c.Addr.(*syscall.SockaddrInet4)
	log.Printf("%d byte read from %d:%d on socket %d\n", sizeMsg, addrFrom.Addr, addrFrom.Port, c.Fd)
	log.Printf("Received command: %s\n", input)

	return commands, nil
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
