package pb

import (
	"testing"

	"github.com/EMCECS/influx/mock"
	ostorage "github.com/influxdata/influxdb/services/storage"
	"google.golang.org/grpc"
	"context"
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
	c := &connection {
		host: hostWithPort,
		conn: cc,
		client: ostorage.NewStorageClient(cc),
	}
	ctx := context.TODO()
	req := &ostorage.ReadRequest{}
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
	c := &connection {
		host: hostWithPort,
		conn: cc,
		client: ostorage.NewStorageClient(cc),
	}
	ctx := context.TODO()
	req := &ostorage.ReadRequest{}

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
