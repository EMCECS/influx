package functions

import (
	"fmt"

	"github.com/influxdata/platform/query"
	"github.com/influxdata/platform/query/execute"
	"github.com/influxdata/platform/query/plan"
	"github.com/influxdata/platform/query/semantic"
	"github.com/influxdata/platform/query/values"
	"github.com/pkg/errors"
)

const PredictLinearKind = "predictLinear"

type PredictLinearOpSpec struct {
	ValueDst           string `json:"value_dst"`
	WantedValue        float64 `json:"wanted_value"`
	execute.AggregateConfig
}

var predictLinearSignature = query.DefaultFunctionSignature()

func init() {
	predictLinearSignature.Params["columns"] = semantic.Array

	query.RegisterBuiltIn("predictLinear", predictLinearBuiltIn)
	query.RegisterFunction(PredictLinearKind, createPredictLinearOpSpec, predictLinearSignature)
	query.RegisterOpSpec(PredictLinearKind, newPredictLinearOp)
	plan.RegisterProcedureSpec(PredictLinearKind, newPredictLinearProcedure, PredictLinearKind)
	execute.RegisterTransformation(PredictLinearKind, createPredictLinearTransformation)
}

// predictLinearBuiltIn defines a `cov` function with an automatic join.
var predictLinearBuiltIn = `
cov = (x,y,on) =>
    join(
        tables:{x:x, y:y},
        on:on,
        fn: (t) => ({x:t.x._value, y:t.y._value}),
    )
    |> predictLinear(columns:["x","y"])

`

func createPredictLinearOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(PredictLinearOpSpec)

	label, ok, err := args.GetString("valueDst")
	if err != nil {
		return nil, err
	} else if ok {
		spec.ValueDst = label
	} else {
		spec.ValueDst = execute.DefaultTimeColLabel
	}

	wantedValue, ok, err := args.GetFloat("wantedValue")
	if err != nil {
		return nil, err
	} else if ok {
		spec.WantedValue = wantedValue
	} else {
		return nil, errors.New("must provide 'wantedValue' argument")
	}

	if err := spec.AggregateConfig.ReadArgs(args); err != nil {
		return nil, err
	}
	if len(spec.Columns) != 2 {
		return nil, errors.New("must provide exactly two columns")
	}
	return spec, nil
}

func newPredictLinearOp() query.OperationSpec {
	return new(PredictLinearOpSpec)
}

func (s *PredictLinearOpSpec) Kind() query.OperationKind {
	return PredictLinearKind
}

type PredictLinearProcedureSpec struct {
	ValueLabel         string
	WantedValue        float64
	execute.AggregateConfig
}

func newPredictLinearProcedure(qs query.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*PredictLinearOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &PredictLinearProcedureSpec{
		ValueLabel:         spec.ValueDst,
		WantedValue:        spec.WantedValue,
		AggregateConfig:    spec.AggregateConfig,
	}, nil
}

func (s *PredictLinearProcedureSpec) Kind() plan.ProcedureKind {
	return PredictLinearKind
}

func (s *PredictLinearProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(PredictLinearProcedureSpec)
	*ns = *s

	ns.AggregateConfig = s.AggregateConfig.Copy()

	return ns
}

type PredictLinearTransformation struct {
	d      execute.Dataset
	cache  execute.TableBuilderCache
	bounds execute.Bounds
	spec   PredictLinearProcedureSpec

	yIdx int

	n,
	symX,
	symY,
	symXY,
	symX2,
	covXY,
	varX,
	slope,
	intercept float64
}

func createPredictLinearTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*PredictLinearProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewPredictLinearTransformation(d, cache, s)
	return t, d, nil
}

func NewPredictLinearTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *PredictLinearProcedureSpec) *PredictLinearTransformation {
	return &PredictLinearTransformation{
		d:     d,
		cache: cache,
		spec:  *spec,
	}
}

func (t *PredictLinearTransformation) RetractTable(id execute.DatasetID, key query.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *PredictLinearTransformation) Process(id execute.DatasetID, tbl query.Table) error {
	cols := tbl.Cols()
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("predictLinear found duplicate table with key: %v", tbl.Key())
	}
	execute.AddTableKeyCols(tbl.Key(), builder)
	valueIdx := builder.AddCol(query.ColMeta{
		Label: t.spec.TimeDst,
		Type:  query.TTime,
	})
	valueIdy := builder.AddCol(query.ColMeta{
		Label: execute.DefaultValueColLabel,
		Type:  query.TFloat,
	})
	yIdx := execute.ColIdx(t.spec.Columns[0], cols)
	xIdx := execute.ColIdx(t.spec.Columns[1], cols)

	// TODO: check if x = float and y = time
	//if cols[yIdx].Type != cols[xIdx].Type {
	//	return errors.New("cannot compute the predictLinear between different types")
	//}
	if err := execute.AppendAggregateTime(t.spec.TimeSrc, t.spec.TimeDst, tbl.Key(), builder); err != nil {
		return err
	}

	t.reset()
	tbl.Do(func(cr query.ColReader) error {
		switch typ := cols[yIdx].Type; typ {
		case query.TFloat:
			t.DoFloat(cr.Floats(yIdx), cr.Times(xIdx))
		default:
			return fmt.Errorf("predictLinear does not support %v", typ)
		}
		return nil
	})

	execute.AppendKeyValues(tbl.Key(), builder)
	builder.AppendTime(valueIdx, t.value())
	builder.AppendFloat(valueIdy, t.spec.WantedValue)
	return nil
}

func (t *PredictLinearTransformation) reset() {
	t.n = 0
	t.symX = 0
	t.symY = 0
	t.symXY = 0
	t.symX2 = 0
}
func (t *PredictLinearTransformation) DoFloat(ys []float64, xs []values.Time) {
	for i, x := range xs {
		y := ys[i]
		x := float64(x)

		t.n++

		t.symX += x
		t.symY += y
		t.symXY += x * y
		t.symX2 += x * x
	}
}
func (t *PredictLinearTransformation) value() values.Time {
	if t.n < 2 {
		return 0
	}
	covXY := t.symXY - t.symX * t.symY/t.n
	varX := t.symX2 - t.symX*t.symX/t.n

	slope := covXY / varX
	intercept := t.symY/t.n - slope*t.symX/t.n

	// predict at which interval value of interest will fall
	predictTime := values.Time((t.spec.WantedValue - intercept ) / slope)

	return predictTime
}

func (t *PredictLinearTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}

func (t *PredictLinearTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}

func (t *PredictLinearTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
