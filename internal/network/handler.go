package network

import (
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/miladbarzideh/goldis/internal/command"
	"github.com/miladbarzideh/goldis/internal/datastore"
	"github.com/miladbarzideh/goldis/utils"
)

const idleTimeout = 60 * time.Second

// ConnectionHandler handles the connection management logic
type ConnectionHandler struct {
	socket         *Socket
	fdMax          int
	activeFd       syscall.FdSet
	fdConn         FdConn
	commandHandler *command.Executor
	idleList       *datastore.DList
	dataStore      *datastore.DataStore
}

// NewConnectionHandler creates a new instance of ConnectionManager
func NewConnectionHandler(socket *Socket) *ConnectionHandler {
	serverFd := socket.Fd
	var activeFd syscall.FdSet
	fdZero(&activeFd)
	fdSet(serverFd, &activeFd)
	fdConn := FdConnInit()
	dataStore := datastore.NewDataStore()
	return &ConnectionHandler{
		socket:         socket,
		fdMax:          serverFd,
		activeFd:       activeFd,
		fdConn:         fdConn,
		commandHandler: command.NewExecutor(dataStore),
		idleList:       datastore.NewDList(),
		dataStore:      dataStore,
	}
}

// StartServer starts the server and handles the connection management logic
func (cm *ConnectionHandler) StartServer() {
	for {
		activeFDSet := cm.getActiveFDSet()
		timeout := cm.nextTimer()
		err := syscall.Select(cm.fdMax+1, &activeFDSet, nil, nil, &timeout)
		if err != nil {
			log.Fatal("Select(): ", err)
		}
		cm.handleActiveConnections(activeFDSet)
		cm.processTimers()
	}
}

func (cm *ConnectionHandler) processTimers() {
	now := time.Now()
	next := cm.idleList.Iterator()
	for nxt := next(); nxt != nil; nxt = next() {
		connection := getConnection(nxt)
		nextTime := connection.idleStart.Add(idleTimeout)
		if nextTime.After(now) {
			break
		}
		addrFrom := connection.Addr.(*syscall.SockaddrInet4)
		log.Printf("Destroy idle connection %d:%d on socket %d\n", addrFrom.Addr, addrFrom.Port, connection.Fd)
		cm.destroyConnection(*connection)
		cm.idleList.Detach(&connection.idleNode, listEq)
	}

	cm.dataStore.RemoveExpiredKeys()
}

func (cm *ConnectionHandler) nextTimer() syscall.Timeval {
	if cm.idleList.IsEmpty() {
		return syscall.Timeval{Sec: 4} // no timer, the value doesn't matter
	}
	now := time.Now()
	connection := getConnection(cm.idleList.GetHead())
	next := connection.idleStart.Add(idleTimeout)
	remaining := next.Sub(now)
	if remaining <= 0 {
		return syscall.Timeval{}
	}
	return syscall.NsecToTimeval(int64(remaining))
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

	cm.resetTimer(connection)

	result := cm.commandHandler.Execute(input)

	_, err = connection.Write([]byte(result + "\n"))
	if err != nil {
		log.Println("Write(): ", err)
	}
}

func (cm *ConnectionHandler) resetTimer(connection Connection) {
	connection.idleStart = time.Now()
	cm.idleList.Detach(&connection.idleNode, listEq)
	cm.idleList.InsertBefore(&connection.idleNode)
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
	cm.idleList.InsertBefore(&connection.idleNode)
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

func listEq(node1, node2 *datastore.LNode) bool {
	c1 := getConnection(node1)
	c2 := getConnection(node2)
	return c1.Fd == c2.Fd
}

func getConnection(node *datastore.LNode) *Connection {
	return (*Connection)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(Connection{}.idleNode)))
}
