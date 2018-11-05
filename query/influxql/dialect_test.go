package influxql_test

import (
	"testing"

	"github.com/EMCECS/influx/query/influxql"
	"github.com/influxdata/flux"
)

func TestDialect(t *testing.T) {
	var _ flux.Dialect = (*influxql.Dialect)(nil)
}
