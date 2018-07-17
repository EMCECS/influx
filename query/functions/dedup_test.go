package functions_test

import (
	"testing"

	"github.com/influxdata/platform/query"
	"github.com/influxdata/platform/query/execute"
	"github.com/influxdata/platform/query/execute/executetest"
	"github.com/influxdata/platform/query/functions"
)

func TestDedup_PassThrough(t *testing.T) {
	executetest.TransformationPassThroughTestHelper(t, func(d execute.Dataset, c execute.BlockBuilderCache) execute.Transformation {
		s := functions.NewDedupTransformation(
			d,
			c,
			&functions.DedupProcedureSpec{
			},
		)
		return s
	})
}

func TestDedup_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *functions.DedupProcedureSpec
		data []query.Block
		want []*executetest.Block
	}{
		{
			name: "one block",
			spec: &functions.DedupProcedureSpec{
				Column: "_value",
			},
			data: []query.Block{&executetest.Block{
				ColMeta: []query.ColMeta{
					{Label: "_time", Type: query.TTime},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(1), 1.0},
					{execute.Time(3), 3.0},
					{execute.Time(4), 1.0},
				},
			}},
			want: []*executetest.Block{{
				ColMeta: []query.ColMeta{
					{Label: "_time", Type: query.TTime},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(3), 3.0},
					{execute.Time(4), 1.0},
				},
			}},
		},
		{
			name: "unique tag",
			spec: &functions.DedupProcedureSpec{
				Column: "t1",
			},
			data: []query.Block{&executetest.Block{
				ColMeta: []query.ColMeta{
					{Label: "_time", Type: query.TTime},
					{Label: "t1", Type: query.TString},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), "a", 2.0},
					{execute.Time(2), "a", 1.0},
					{execute.Time(3), "b", 3.0},
					{execute.Time(2), "a", 1.0},
				},
			}},
			want: []*executetest.Block{{
				ColMeta: []query.ColMeta{
					{Label: "_time", Type: query.TTime},
					{Label: "t1", Type: query.TString},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), "a", 2.0},
					{execute.Time(2), "a", 1.0},
					{execute.Time(3), "b", 3.0},
				},
			}},
		},
		{
			name: "unique times",
			spec: &functions.DedupProcedureSpec{
				Column: "_time",
			},
			data: []query.Block{&executetest.Block{
				ColMeta: []query.ColMeta{
					{Label: "_time", Type: query.TTime},
					{Label: "t1", Type: query.TString},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), "a", 2.0},
					{execute.Time(2), "a", 1.0},
					{execute.Time(3), "b", 3.0},
					{execute.Time(3), "c", 1.0},
				},
			}},
			want: []*executetest.Block{{
				ColMeta: []query.ColMeta{
					{Label: "_time", Type: query.TTime},
					{Label: "t1", Type: query.TString},
					{Label: "_value", Type: query.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), "a", 2.0},
					{execute.Time(2), "a", 1.0},
					{execute.Time(3), "b", 3.0},
					{execute.Time(3), "c", 1.0},
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
				func(d execute.Dataset, c execute.BlockBuilderCache) execute.Transformation {
					return functions.NewDedupTransformation(d, c, tc.spec)
				},
			)
		})
	}
}
