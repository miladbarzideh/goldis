package network

import (
	"net"
	"syscall"
)

type Socket struct {
	fd int
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
	return &Socket{fd: fd}, nil
}

func (s Socket) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n, err := syscall.Read(s.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}

func (s Socket) Write(p []byte) (int, error) {
	n, err := syscall.Write(s.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}

func (s Socket) Accept() (*Socket, error) {
	nfd, _, err := syscall.Accept(s.fd)
	if err != nil {
		return nil, err
	}
	syscall.CloseOnExec(nfd)
	return &Socket{nfd}, nil
}

func (s Socket) Close() error {
	return syscall.Close(s.fd)
}
