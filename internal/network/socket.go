package network

import (
	"net"
	"syscall"
)

type Socket struct {
	Fd int
}

func NewSocket(ip net.IP, port int) (*Socket, error) {
	//AF_INET for IPv4 & SOCK_STREAM for TCP connection
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	//configure the socket
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return nil, err
	}

	//bind on a wildcard address 0.0.0.0:8585
	socketAddress := &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte(ip),
	}
	if err := syscall.Bind(fd, socketAddress); err != nil {
		return nil, err
	}

	//listen
	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return nil, err
	}
	return &Socket{Fd: fd}, nil
}

func (s Socket) Accept() (*Connection, error) {
	fd, addr, err := syscall.Accept(s.Fd)
	if err != nil {
		return nil, err
	}
	syscall.CloseOnExec(fd)
	return &Connection{Fd: fd, Addr: addr}, nil
}

func (s Socket) Close() error {
	return syscall.Close(s.Fd)
}
