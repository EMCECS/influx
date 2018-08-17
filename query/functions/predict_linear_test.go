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

func TestPredictLinear_NewQuery(t *testing.T) {
	tests := []querytest.NewQueryTestCase{
		{
			Name: "simple regression",
			Raw:  `from(db:"mydb") |> predictLinear(columns:["a","b"], wantedValue: 10.0)`,
			Want: &query.Spec{
				Operations: []*query.Operation{
					{
						ID: "from0",
						Spec: &functions.FromOpSpec{
							Database: "mydb",
						},
					},
					{
						ID: "predictLinear1",
						Spec: &functions.PredictLinearOpSpec{
							ValueDst: execute.DefaultTimeColLabel,
							WantedValue: 10.0,
							AggregateConfig: execute.AggregateConfig{
								TimeSrc: execute.DefaultStopColLabel,
								TimeDst: execute.DefaultTimeColLabel,
								Columns: []string{"a", "b"},
							},
						},
					},
				},
				Edges: []query.Edge{
					{Parent: "from0", Child: "predictLinear1"},
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			querytest.NewQueryTestHelper(t, tc)
		})
	}
}

func TestPredictLinear_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *functions.PredictLinearProcedureSpec
		data []query.Table
		want []*executetest.Table
	}{
		{
			name: "simple regression",
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
		{
			name: "earlier time",
			spec: &functions.PredictLinearProcedureSpec{
				WantedValue: 0,
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
					{execute.Time(0), execute.Time(5), execute.Time(10), 1.0},
					{execute.Time(0), execute.Time(5), execute.Time(11), 2.0},
					{execute.Time(0), execute.Time(5), execute.Time(12), 3.0},
					{execute.Time(0), execute.Time(5), execute.Time(13), 4.0},
					{execute.Time(0), execute.Time(5), execute.Time(14), 5.0},
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
					{execute.Time(0), execute.Time(5), execute.Time(9), 0.0},
				},
			}},
		},
		{
			name: "negative time",
			spec: &functions.PredictLinearProcedureSpec{
				WantedValue: 0,
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
					{execute.Time(0), execute.Time(5), execute.Time(-1), 0.0},
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
