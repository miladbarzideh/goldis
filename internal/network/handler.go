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
	activeFd       syscall.FdSet
	fdConn         FdConn
	commandHandler *command.Handler
}

// NewConnectionHandler creates a new instance of ConnectionManager
func NewConnectionHandler(socket *Socket) *ConnectionHandler {
	serverFd := socket.Fd
	var activeFd syscall.FdSet
	fdZero(&activeFd)
	fdSet(serverFd, &activeFd)
	fdConn := FdConnInit()
	return &ConnectionHandler{
		socket:         socket,
		fdMax:          serverFd,
		activeFd:       activeFd,
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
	fdSet(acceptedFd, &cm.activeFd)
	cm.fdConn.set(acceptedFd, *connection)
	if acceptedFd > cm.fdMax {
		cm.fdMax = acceptedFd
	}
}

func (cm *ConnectionHandler) destroyConnection(connection Connection) {
	fdClr(connection.Fd, &cm.activeFd)
	cm.fdConn.clr(connection.Fd)
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
	tmpFDSet := cm.activeFd
	return tmpFDSet
}

func (cm *ConnectionHandler) handleActiveConnections(fdSet syscall.FdSet) {
	// the event loop
	for fd := 0; fd < cm.fdMax+1; fd++ {
		if fdIsSet(fd, &fdSet) {
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

func fdIsSet(fd int, p *syscall.FdSet) bool {
	return p.Bits[fd/32]&(1<<(uint(fd)%32)) != 0
}

func fdSet(fd int, p *syscall.FdSet) {
	p.Bits[fd/32] |= 1 << (uint(fd) % 32)
}

func fdClr(fd int, p *syscall.FdSet) {
	p.Bits[fd/32] &^= 1 << (uint(fd) % 32)
}

type FdConn map[int]Connection

func FdConnInit() FdConn {
	return make(FdConn, syscall.FD_SETSIZE)
}

func (f *FdConn) set(fd int, value Connection) {
	(*f)[fd] = value
}

func (f *FdConn) clr(fd int) {
	delete(*f, fd)
}

func fdZero(p *syscall.FdSet) {
	p.Bits = [32]int32{}
}
