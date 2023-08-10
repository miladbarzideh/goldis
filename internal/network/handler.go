package network

import (
	"log"
	"syscall"

	"github.com/miladbarzideh/goldis/internal/command"
)

// ConnectionHandler handles the connection management logic
type ConnectionHandler struct {
	socket         *Socket
	fdMax          int
	activeFD       syscall.FdSet
	fdConn         FdConn
	commandHandler *command.Handler
}

// NewConnectionHandler creates a new instance of ConnectionManager
func NewConnectionHandler(socket *Socket) *ConnectionHandler {
	serverFd := socket.Fd
	var activeFd syscall.FdSet
	FDZero(&activeFd)
	FDSet(serverFd, &activeFd)
	fdConn := FDConnInit()
	return &ConnectionHandler{
		socket:         socket,
		fdMax:          serverFd,
		activeFD:       activeFd,
		fdConn:         fdConn,
		commandHandler: command.NewHandler(),
	}
}

// StartServer starts the server and handles the connection management logic
func (cm *ConnectionHandler) StartServer() {
	for {
		activeFDSet := cm.getActiveFDSet()
		err := syscall.Select(cm.fdMax+1, &activeFDSet, nil, nil, nil)
		if err != nil {
			log.Fatal("Select(): ", err)
		}
		cm.handleActiveConnections(activeFDSet)
	}
}

func (cm *ConnectionHandler) addConnection(connection *Connection) {
	acceptedFd := connection.Fd
	FDSet(acceptedFd, &cm.activeFD)
	cm.fdConn.Set(acceptedFd, *connection)
	if acceptedFd > cm.fdMax {
		cm.fdMax = acceptedFd
	}
}

func (cm *ConnectionHandler) destroyConnection(connection Connection) {
	FDClr(connection.Fd, &cm.activeFD)
	cm.fdConn.Clr(connection.Fd)
	_ = connection.Close()
}

func (cm *ConnectionHandler) handleConnectionIO(connection Connection) {
	input, err := connection.Read()
	if err != nil {
		log.Println("Read(): ", err)
		cm.destroyConnection(connection)
		return
	}

	result := cm.commandHandler.Execute(input)

	_, err = connection.Write([]byte(result + "\n"))
	if err != nil {
		log.Println("Write(): ", err)
	}
}

func (cm *ConnectionHandler) getActiveFDSet() syscall.FdSet {
	tmpFDSet := cm.activeFD
	return tmpFDSet
}

func (cm *ConnectionHandler) handleActiveConnections(fdSet syscall.FdSet) {
	// the event loop
	for fd := 0; fd < cm.fdMax+1; fd++ {
		if FDIsSet(fd, &fdSet) {
			if fd == cm.socket.Fd {
				cm.acceptNewConnection()
			} else {
				cm.handleConnectionIO(cm.fdConn[fd])
			}
		}
	}
}

func (cm *ConnectionHandler) acceptNewConnection() {
	connection, err := cm.socket.Accept()
	if err != nil {
		log.Fatal("Accept(): ", err)
	}
	cm.addConnection(connection)
}

func FDIsSet(fd int, p *syscall.FdSet) bool {
	return p.Bits[fd/32]&(1<<(uint(fd)%32)) != 0
}

func FDSet(fd int, p *syscall.FdSet) {
	p.Bits[fd/32] |= 1 << (uint(fd) % 32)
}

func FDClr(fd int, p *syscall.FdSet) {
	p.Bits[fd/32] &^= 1 << (uint(fd) % 32)
}

type FdConn map[int]Connection

func FDConnInit() FdConn {
	return make(FdConn, syscall.FD_SETSIZE)
}

func (f *FdConn) Get(fd int) Connection {
	return (*f)[fd]
}

func (f *FdConn) Set(fd int, value Connection) {
	(*f)[fd] = value
}

func (f *FdConn) Clr(fd int) {
	delete(*f, fd)
}

func FDZero(p *syscall.FdSet) {
	p.Bits = [32]int32{}
}
