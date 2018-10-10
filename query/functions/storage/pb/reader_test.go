package pb

import (
	"testing"

	"context"
	"github.com/akurilov/influx/mock"
	"google.golang.org/grpc"
	"strconv"
)

const (
	DefaultPort = 8082
)

func TestReader_Success(t *testing.T) {
	server := mock.NewServerMock(mock.DefaultConnProt, mock.DefaultConnHost, DefaultPort)
	defer server.Close()
	err := server.Start()
	if err != nil {
		t.Error(err)
	}
	hostWithPort := mock.DefaultConnHost + ":" + strconv.Itoa(DefaultPort)
	cc, err := grpc.Dial(hostWithPort, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	c := &connection{
		host:   hostWithPort,
		conn:   cc,
		client: NewStorageClient(cc),
	}
	ctx := context.TODO()
	req := &ReadRequest{}
	stream, newConn, err := readWithRecovery(c, &ctx, req)
	if stream == nil {
		t.Error("no stream returned")
	}
	if newConn != nil {
		t.Error("unexpected reconnect")
	}
	if err != nil {
		t.Error(err)
	}

}

func TestReader_ConnRecovery(t *testing.T) {

	var server = mock.NewServerMock(mock.DefaultConnProt, mock.DefaultConnHost, DefaultPort)
	var err = server.Start()
	if err != nil {
		t.Error(err)
	}

	hostWithPort := mock.DefaultConnHost + ":" + strconv.Itoa(DefaultPort)
	cc, err := grpc.Dial(hostWithPort, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	c := &connection{
		host:   hostWithPort,
		conn:   cc,
		client: NewStorageClient(cc),
	}
	ctx := context.TODO()
	req := &ReadRequest{}

	server.Close()
	server = mock.NewServerMock(mock.DefaultConnProt, mock.DefaultConnHost, DefaultPort)
	err = server.Start()
	if err != nil {
		t.Error(err)
	}

	stream, _, err := readWithRecovery(c, &ctx, req)
	if stream == nil {
		t.Error("no stream returned")
	}
	if err != nil {
		t.Error(err)
	}
}
