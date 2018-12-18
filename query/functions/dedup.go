package functions

import (
	"fmt"

	"github.com/EMCECS/influx/query"
	"github.com/EMCECS/influx/query/execute"
	"github.com/EMCECS/influx/query/plan"
	"github.com/EMCECS/influx/query/semantic"
)

const DedupKind = "dedup"

type DedupOpSpec struct {
}

var dedupSignature = query.DefaultFunctionSignature()

func init() {
	dedupSignature.Params["column"] = semantic.String

	query.RegisterFunction(DedupKind, createDedupOpSpec, dedupSignature)
	query.RegisterOpSpec(DedupKind, newDedupOp)
	plan.RegisterProcedureSpec(DedupKind, newDedupProcedure, DedupKind)
	execute.RegisterTransformation(DedupKind, createDedupTransformation)
}

func createDedupOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(DedupOpSpec)

	return spec, nil
}

func newDedupOp() query.OperationSpec {
	return new(DedupOpSpec)
}

func (s *DedupOpSpec) Kind() query.OperationKind {
	return DedupKind
}

type DedupProcedureSpec struct {
	Column string
}

func newDedupProcedure(qs query.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	_, ok := qs.(*DedupOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &DedupProcedureSpec{
	}, nil
}

func (s *DedupProcedureSpec) Kind() plan.ProcedureKind {
	return DedupKind
}
func (s *DedupProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(DedupProcedureSpec)

	*ns = *s

	return ns
}

func createDedupTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*DedupProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewDedupTransformation(d, cache, s)
	return t, d, nil
}

type dedupTransformation struct {
	d     execute.Dataset
	cache execute.TableBuilderCache
}

func NewDedupTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *DedupProcedureSpec) *dedupTransformation {
	return &dedupTransformation{
		d:      d,
		cache:  cache,
	}
}

func (t *dedupTransformation) RetractTable(id execute.DatasetID, key query.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *dedupTransformation) Process(id execute.DatasetID, tbl query.Table) error {

	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("dedup found duplicate block with key: %v", tbl.Key())
	}
	execute.AddTableCols(tbl, builder)

	var (
		lastBool   map[int]bool
		lastInt    map[int]int64
		lastUint   map[int]uint64
		lastFloat  map[int]float64
		lastString map[int]string
		lastTime   map[int]execute.Time
	)
	lastBool = make(map[int]bool)
	lastInt = make(map[int]int64)
	lastUint = make(map[int]uint64)
	lastFloat = make(map[int]float64)
	lastString = make(map[int]string)
	lastTime = make(map[int]execute.Time)

	firstRow := true

	return tbl.Do(func(cr query.ColReader) error {

		l := cr.Len()
		builderCols := builder.Cols()
		colCount := len(builderCols)

		// loop over the records
		for i := 0; i < l; i ++ {

			duplicateFlag := !firstRow // assume the row to be a duplicate if the row is not 1st

			for j := 0; j < colCount; j ++ { // loop over the columns
				col := builderCols[j]
				switch col.Type {
				case query.TBool:
					if cr.Bools(j)[i] != lastBool[j] {
						duplicateFlag = false
					}
					lastBool[j] = cr.Bools(j)[i]
				case query.TInt:
					if cr.Ints(j)[i] != lastInt[j] {
						duplicateFlag = false
					}
					lastInt[j] = cr.Ints(j)[i]
				case query.TUInt:
					if cr.UInts(j)[i] != lastUint[j] {
						duplicateFlag = false
					}
					lastUint[j] = cr.UInts(j)[i]
				case query.TFloat:
					if cr.Floats(j)[i] != lastFloat[j] {
						duplicateFlag = false
					}
					lastFloat[j] = cr.Floats(j)[i]
				case query.TString:
					if cr.Strings(j)[i] != lastString[j] {
						duplicateFlag = false
					}
					lastString[j] = cr.Strings(j)[i]
				case query.TTime:
					if cr.Times(j)[i] != lastTime[j] {
						duplicateFlag = false
					}
					lastTime[j] = cr.Times(j)[i]
				}
			}

			if !duplicateFlag {
				execute.AppendRecord(i, cr, builder)
			}

			if firstRow {
				firstRow = false
			}
		}
		return nil
	})
}

func (t *dedupTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}
func (t *dedupTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}
func (t *dedupTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
