package influxql_test

import (
	"testing"

	"github.com/EMCECS/influx/query"
	"github.com/EMCECS/influx/query/influxql"
)

func TestCompiler(t *testing.T) {
	var _ query.Compiler = (*influxql.Compiler)(nil)
}
