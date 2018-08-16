package functions_test

import (
	"testing"

	"github.com/influxdata/platform/query"
	"github.com/influxdata/platform/query/execute/executetest"
	"github.com/influxdata/platform/query/functions"
	"github.com/influxdata/platform/query/querytest"
	"github.com/influxdata/platform/query/execute"
)

func TestPredictLinearOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"predictLinear","kind":"predictLinear"}`)
	op := &query.Operation{
		ID:   "predictLinear",
		Spec: &functions.PredictLinearOpSpec{},
	}

	querytest.OperationMarshalingTestHelper(t, data, op)
}


func TestPredictLinear_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *functions.PredictLinearProcedureSpec
		data []query.Table
		want []*executetest.Table
	}{
		{
			name: "variance",
			spec: &functions.PredictLinearProcedureSpec{
				WantedValue: 50,
				ValueLabel: execute.DefaultTimeColLabel,
				AggregateConfig: execute.AggregateConfig{
					TimeSrc: execute.DefaultStopColLabel,
					TimeDst: execute.DefaultTimeColLabel,
					Columns: []string{"x", "_time"},
				},
			},
			data: []query.Table{&executetest.Table{
				KeyCols: []string{"_start", "_stop"},
				ColMeta: []query.ColMeta{
					{Label: "_start", Type: query.TTime},
					{Label: "_stop", Type: query.TTime},
					{Label: "_time", Type: query.TTime},
					{Label: "x", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(0), execute.Time(5), execute.Time(0), 1.0},
					{execute.Time(0), execute.Time(5), execute.Time(1), 2.0},
					{execute.Time(0), execute.Time(5), execute.Time(2), 3.0},
					{execute.Time(0), execute.Time(5), execute.Time(3), 4.0},
					{execute.Time(0), execute.Time(5), execute.Time(4), 5.0},
				},
			}},
			want: []*executetest.Table{{
				KeyCols: []string{"_start", "_stop"},
				ColMeta: []query.ColMeta{
					{Label: "_start", Type: query.TTime},
					{Label: "_stop", Type: query.TTime},
					{Label: "_time", Type: query.TTime},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(0), execute.Time(5), execute.Time(49), 50.0},
				},
			}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.ProcessTestHelper(
				t,
				tc.data,
				tc.want,
				func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
					return functions.NewPredictLinearTransformation(d, c, tc.spec)
				},
			)
		})
	}
}

func BenchmarkPredictLinear(b *testing.B) {
	executetest.AggFuncBenchmarkHelper(
		b,
		new(functions.PredictLinearTransformation),
		NormalData,
		10.00081696729983,
	)
}
