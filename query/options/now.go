package options

import (
	"github.com/EMCECS/influx/query"
	"github.com/EMCECS/influx/query/functions"
)

func init() {
	query.RegisterBuiltInOption("now", functions.SystemTime())
}
