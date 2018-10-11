package executetest

import (
	"math"

	"github.com/EMCECS/influx/query/execute"
)

var UnlimitedAllocator = &execute.Allocator{
	Limit: math.MaxInt64,
}
