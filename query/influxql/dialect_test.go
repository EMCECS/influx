package influxql_test

import (
	"testing"

	"github.com/EMCECS/influx/query"
	"github.com/EMCECS/influx/query/influxql"
)

func TestDialect(t *testing.T) {
	var _ query.Dialect = (*influxql.Dialect)(nil)
}
