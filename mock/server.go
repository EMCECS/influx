package mock

import(
	"net"
	"fmt"
	"strconv"
)

const(
	DefaultConnProt = "tcp"
	DefaultConnHost = "localhost"
	DefaultBuffSize = 1024
)

type ServerMock struct {
	prot		string
	host 		string
	port 		uint16
	listener	net.Listener
	conns		[]net.Conn
}

func NewServerMock(prot string, host string, port uint16) (*ServerMock) {
	return &ServerMock{
		prot: 		prot,
		host:		host,
		port:		port,
		listener: 	nil,
		conns:		make([]net.Conn, 1024),
	}
}

func (this *ServerMock) Start() (error) {
	listener, err := net.Listen(this.prot, this.host+":"+strconv.Itoa(int(this.port)))
	this.listener = listener
	return err
}

func (this *ServerMock) Close() (error) {
	var lastErr error
	err := this.listener.Close()
	if err != nil {
		lastErr = err
	}
	for _, conn := range this.conns {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				lastErr = err
			}
		}
	}
	return lastErr
}

func (this *ServerMock) acceptLoop(listener net.Listener) {
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err == nil {
			this.conns = append(this.conns, conn)
			// Handle connections in a new goroutine.
			go handleRequest(conn)
		}
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, DefaultBuffSize)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("TCP server mock error reading:", err.Error())
	}
	conn.Write(buf)
}
