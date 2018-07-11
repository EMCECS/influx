package pb

import (
	"testing"

	"github.com/influxdata/platform/mock"
	"google.golang.org/grpc"
)

const (
	DefaultPort = 8082
)

func TestReader_ConnRecovery(t *testing.T)  {
	var serverMock, _ = mock.NewServerMock(mock.DefaultConnProt, mock.DefaultConnHost, DefaultPort)
	go serverMock.Listen()
	cc, err := grpc.Dial(mock.DefaultConnHost + ":" + string(DefaultPort))
	stream, newConn, err = readWithRecovery(c, ctx, req)
}

func TestReader_ConnShutdown(t *testing.T) {

}