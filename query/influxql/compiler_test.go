package influxql_test

import (
	"testing"

	"github.com/EMCECS/influx/query/influxql"
	"github.com/influxdata/flux"
)

func TestCompiler(t *testing.T) {
	var _ flux.Compiler = (*influxql.Compiler)(nil)
}
