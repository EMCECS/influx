package reads

import (
	"context"

	"github.com/EMCECS/influx/models"
	fstorage "github.com/EMCECS/influx/query/functions/inputs/storage"
	"github.com/EMCECS/influx/storage/reads/datatypes"
	"github.com/EMCECS/influx/tsdb/cursors"
	"github.com/gogo/protobuf/proto"
)

type ResultSet interface {
	Close()
	Next() bool
	Cursor() cursors.Cursor
	Tags() models.Tags
}

type GroupResultSet interface {
	Next() GroupCursor
	Close()
}

type GroupCursor interface {
	Tags() models.Tags
	Keys() [][]byte
	PartitionKeyVals() [][]byte
	Next() bool
	Cursor() cursors.Cursor
	Close()
}

type Store interface {
	Read(ctx context.Context, req *datatypes.ReadRequest) (ResultSet, error)
	GroupRead(ctx context.Context, req *datatypes.ReadRequest) (GroupResultSet, error)
	GetSource(rs fstorage.ReadSpec) (proto.Message, error)
}
