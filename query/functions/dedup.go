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

func (t *dedupTransformation) Process(id execute.DatasetID, b query.Table) error {

	builder, created := t.cache.TableBuilder(b.Key())
	if !created {
		return fmt.Errorf("dedup found duplicate block with key: %v", b.Key())
	}
	execute.AddTableCols(b, builder)

	var (
		boolDedup   map[bool]bool
		intDedup    map[int64]bool
		uintDedup   map[uint64]bool
		floatDedup  map[float64]bool
		stringDedup map[string]bool
		timeDedup   map[execute.Time]bool
	)
	boolDedup = make(map[bool]bool)
	intDedup = make(map[int64]bool)
	uintDedup = make(map[uint64]bool)
	floatDedup = make(map[float64]bool)
	stringDedup = make(map[string]bool)
	timeDedup = make(map[execute.Time]bool)

	return b.Do(func(cr query.ColReader) error {
		l := cr.Len()
		colCount := len(builder.Cols())
		// loop over the records
		for i := 0; i < l; i ++ {
			duplicateFlag := true
			// loop over the columns
			for j := 0; j < colCount; j ++ {
				col := builder.Cols()[j]
				switch col.Type {
				case query.TBool:
					v := cr.Bools(j)[i]
					if boolDedup[v] {
						continue
					} else {
						duplicateFlag = false
					}
					boolDedup[v] = true
				case query.TInt:
					v := cr.Ints(j)[i]
					if intDedup[v] {
						continue
					} else {
						duplicateFlag = false
					}
					intDedup[v] = true
				case query.TUInt:
					v := cr.UInts(j)[i]
					if uintDedup[v] {
						continue
					} else {
						duplicateFlag = false
					}
					uintDedup[v] = true
				case query.TFloat:
					v := cr.Floats(j)[i]
					if floatDedup[v] {
						continue
					} else {
						duplicateFlag = false
					}
					floatDedup[v] = true
				case query.TString:
					v := cr.Strings(j)[i]
					if stringDedup[v] {
						continue
					} else {
						duplicateFlag = false
					}
					stringDedup[v] = true
				case query.TTime:
					v := cr.Times(j)[i]
					if timeDedup[v] {
						continue
					} else {
						duplicateFlag = false
					}
					timeDedup[v] = true
				}
			}

			if !duplicateFlag {
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
