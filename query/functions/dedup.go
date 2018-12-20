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

	return &DedupProcedureSpec{}, nil
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
		d:     d,
		cache: cache,
	}
}

func (t *dedupTransformation) RetractTable(id execute.DatasetID, key query.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *dedupTransformation) Process(id execute.DatasetID, tbl query.Table) error {

	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("dedup found duplicate table with key: %v", tbl.Key())
	}
	execute.AddTableCols(tbl, builder)
	builderCols := builder.Cols()
	colCount := len(builderCols)
	colIdxTimestamp := -1
	for i := 0; i < colCount; i++ {
		col := builderCols[i]
		if "_time" == col.Label {
			colIdxTimestamp = i
		}
	}
	if colIdxTimestamp < 0 {
		return fmt.Errorf("dedup: column _time not found in the table with key %v", tbl.Key())
	}
	uniqueTimestamps := make(map[execute.Time]struct{})

	return tbl.Do(func(cr query.ColReader) error {
		l := cr.Len()
		// loop over the records
		for i := 0; i < l; i++ {
			ts := cr.Times(colIdxTimestamp)[i]
			_, duplicateFlag := uniqueTimestamps[ts]
			if !duplicateFlag {
				uniqueTimestamps[ts] = struct{}{}
				execute.AppendRecord(i, cr, builder)
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
