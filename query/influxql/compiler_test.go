package influxql_test

import (
	"testing"

	"github.com/EMCECS/flux"
	"github.com/EMCECS/influx/query/influxql"
)

func TestCompiler(t *testing.T) {
	var _ flux.Compiler = (*influxql.Compiler)(nil)
}
