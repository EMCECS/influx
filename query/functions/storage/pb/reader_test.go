package pb

import (
	"strconv"
	"testing"
	"time"

	"github.com/EMCECS/influx/mock"
	"google.golang.org/grpc"
)

const (
	DefaultPort = 8082
)

func TestWaitConnForReady(t *testing.T) {
	//setup mock server
	server := mock.NewServerMock(mock.DefaultConnProt, mock.DefaultConnHost, DefaultPort)
	defer server.Close()
	err := server.Start()
	if err != nil {
		t.Error(err)
	}
	hostWithPort := mock.DefaultConnHost + ":" + strconv.Itoa(DefaultPort)
	//setup grpc connections with server
	cc, err := grpc.Dial(hostWithPort, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer cc.Close()
	cc2, err := grpc.Dial(hostWithPort, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer cc2.Close()
	cc3, err := grpc.Dial(hostWithPort, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer cc3.Close()
	//setup connection with non exist host that will always fail to connect
	cctimeout, err := grpc.Dial("passthrough:///Non-Existent.Server:80", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer cctimeout.Close()

	type args struct {
		conns []connection
		wait  time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Ready multiple conn",
			args{
				[]connection{
					connection{conn: cc},
					connection{conn: cc2},
					connection{conn: cc3},
				},
				10 * time.Millisecond,
			},
			false,
		},
		{"Ready one",
			args{
				[]connection{
					connection{conn: cc},
					connection{conn: cc2},
					connection{conn: cc3},
				},
				10 * time.Millisecond,
			},
			false,
		},
		{"Cancelled",
			args{
				[]connection{
					connection{conn: cc},
					connection{conn: cc2},
					connection{conn: cc3},
					connection{conn: cctimeout},
				},
				10 * time.Millisecond,
			},
			true,
		},
		{"Empty",
			args{
				[]connection{},
				10 * time.Millisecond,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WaitConnForReady(tt.args.conns, tt.args.wait); (err != nil) != tt.wantErr {
				t.Errorf("WaitConnForReady() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
