package mock

import (
	"net"
	"fmt"
	"os"
)

const (
	DefaultConnProt = "tcp"
	DefaultConnHost = "localhost"
	DefaultBuffSize = 1024
)

type ServerMock struct {
	prot	string
	host    string
	port	uint16
}

func NewServerMock(prot string, host string, port uint16) (*ServerMock, error) {
	return &ServerMock {
		prot: prot,
		host: host,
		port: port,
	}, nil
}

func (s *ServerMock) Listen() {

	l, err := net.Listen(s.prot, s.host + ":" + string(s.port))
	if err != nil {
		fmt.Println("TCP server mock error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("TCP server mock is listening on " + s.host + ":" + string(s.port))
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("TCP server mock error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
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
	// Close the connection when you're done with it.
	conn.Close()
}
