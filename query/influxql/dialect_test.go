package influxql_test

import (
	"testing"

	"github.com/EMCECS/flux"
	"github.com/EMCECS/influx/query/influxql"
)

func TestDialect(t *testing.T) {
	var _ flux.Dialect = (*influxql.Dialect)(nil)
}
