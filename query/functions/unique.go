package functions

import (
	"fmt"

	"github.com/influxdata/platform/query"
	"github.com/influxdata/platform/query/execute"
	"github.com/influxdata/platform/query/plan"
	"github.com/influxdata/platform/query/semantic"
)

const UniqueKind = "unique"

type UniqueOpSpec struct {
	Column string `json:"column"`
	All bool `json:"all"`
}

var uniqueSignature = query.DefaultFunctionSignature()

func init() {
	uniqueSignature.Params["column"] = semantic.String
	uniqueSignature.Params["all"] = semantic.Bool

	query.RegisterFunction(UniqueKind, createUniqueOpSpec, uniqueSignature)
	query.RegisterOpSpec(UniqueKind, newUniqueOp)
	plan.RegisterProcedureSpec(UniqueKind, newUniqueProcedure, UniqueKind)
	execute.RegisterTransformation(UniqueKind, createUniqueTransformation)
}

func createUniqueOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(UniqueOpSpec)

	if col, ok, err := args.GetString("column"); err != nil {
		return nil, err
	} else if ok {
		spec.Column = col
	} else {
		spec.Column = execute.DefaultValueColLabel
	}

	if all, ok, err := args.GetBool("all"); err != nil {
		return nil, err
	} else if ok {
		spec.All = all
	} else {
		spec.All = false
	}

	return spec, nil
}

func newUniqueOp() query.OperationSpec {
	return new(UniqueOpSpec)
}

func (s *UniqueOpSpec) Kind() query.OperationKind {
	return UniqueKind
}

type UniqueProcedureSpec struct {
	Column string
	All bool
}

func newUniqueProcedure(qs query.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*UniqueOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &UniqueProcedureSpec{
		Column: spec.Column,
		All:	spec.All,
	}, nil
}

func (s *UniqueProcedureSpec) Kind() plan.ProcedureKind {
	return UniqueKind
}
func (s *UniqueProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(UniqueProcedureSpec)

	*ns = *s

	return ns
}

func createUniqueTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*UniqueProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewUniqueTransformation(d, cache, s)
	return t, d, nil
}

type uniqueTransformation struct {
	d     execute.Dataset
	cache execute.TableBuilderCache

	column	string
	all		bool
}

func NewUniqueTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *UniqueProcedureSpec) *uniqueTransformation {
	return &uniqueTransformation{
		d:      d,
		cache:  cache,
		column: spec.Column,
		all:	spec.All,
	}
}

func (t *uniqueTransformation) RetractTable(id execute.DatasetID, key query.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *uniqueTransformation) Process(id execute.DatasetID, tbl query.Table) error {
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("unique found duplicate table with key: %v", tbl.Key())
	}
	execute.AddTableCols(tbl, builder)

	var (
		boolUnique   map[bool]bool
		intUnique    map[int64]bool
		uintUnique   map[uint64]bool
		floatUnique  map[float64]bool
		stringUnique map[string]bool
		timeUnique   map[execute.Time]bool
	)

	if t.all {

		boolUnique= make(map[bool]bool)
		intUnique = make(map[int64]bool)
		uintUnique = make(map[uint64]bool)
		floatUnique = make(map[float64]bool)
		stringUnique = make(map[string]bool)
		timeUnique = make(map[execute.Time]bool)

		return tbl.Do(func(cr query.ColReader) error {
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
						if boolUnique[v] {
							continue
						} else {
							duplicateFlag = false
						}
						boolUnique[v] = true
					case query.TInt:
						v := cr.Ints(j)[i]
						if intUnique[v] {
							continue
						} else {
							duplicateFlag = false
						}
						intUnique[v] = true
					case query.TUInt:
						v := cr.UInts(j)[i]
						if uintUnique[v] {
							continue
						} else {
							duplicateFlag = false
						}
						uintUnique[v] = true
					case query.TFloat:
						v := cr.Floats(j)[i]
						if floatUnique[v] {
							continue
						} else {
							duplicateFlag = false
						}
						floatUnique[v] = true
					case query.TString:
						v := cr.Strings(j)[i]
						if stringUnique[v] {
							continue
						} else {
							duplicateFlag = false
						}
						stringUnique[v] = true
					case query.TTime:
						v := cr.Times(j)[i]
						if timeUnique[v] {
							continue
						} else {
							duplicateFlag = false
						}
						timeUnique[v] = true
					}
				}

				if !duplicateFlag {
					execute.AppendRecord(i, cr, builder)
				}
			}
			return nil
		})

	} else {

		colIdx := execute.ColIdx(t.column, builder.Cols())
		if colIdx < 0 {
			return fmt.Errorf("no column %q exists", t.column)
		}
		col := builder.Cols()[colIdx]

		switch col.Type {
		case query.TBool:
			boolUnique = make(map[bool]bool)
		case query.TInt:
			intUnique = make(map[int64]bool)
		case query.TUInt:
			uintUnique = make(map[uint64]bool)
		case query.TFloat:
			floatUnique = make(map[float64]bool)
		case query.TString:
			stringUnique = make(map[string]bool)
		case query.TTime:
			timeUnique = make(map[execute.Time]bool)
		}

	return tbl.Do(func(cr query.ColReader) error {
		l := cr.Len()
		for i := 0; i < l; i++ {
			// Check unique
			switch col.Type {
			case query.TBool:
				v := cr.Bools(colIdx)[i]
				if boolUnique[v] {
					continue
				}
				boolUnique[v] = true
			case query.TInt:
				v := cr.Ints(colIdx)[i]
				if intUnique[v] {
					continue
				}
				intUnique[v] = true
			case query.TUInt:
				v := cr.UInts(colIdx)[i]
				if uintUnique[v] {
					continue
				}
				uintUnique[v] = true
			case query.TFloat:
				v := cr.Floats(colIdx)[i]
				if floatUnique[v] {
					continue
				}
				floatUnique[v] = true
			case query.TString:
				v := cr.Strings(colIdx)[i]
				if stringUnique[v] {
					continue
				}
				stringUnique[v] = true
			case query.TTime:
				v := cr.Times(colIdx)[i]
				if timeUnique[v] {
					continue
				}
				timeUnique[v] = true
			}

				execute.AppendRecord(i, cr, builder)
			}
			return nil
		})
	}
}

func (t *uniqueTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}
func (t *uniqueTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}
func (t *uniqueTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
