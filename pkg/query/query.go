package query

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

type Query interface {
	Query(field string, builder squirrel.SelectBuilder) (squirrel.SelectBuilder, error)
}

var ErrInvalidFilterOperator = errors.New("Invalid Filter Operation")

type SortBy struct {
	Field string
	Asc   bool
}

type SortBys []SortBy

func (s SortBys) Query(field string, builder squirrel.SelectBuilder) (squirrel.SelectBuilder, error) {
	if len(s) == 0 {
		return builder, nil
	}

	for _, sortBy := range s {
		field := sortBy.Field
		if !sortBy.Asc {
			field += " DESC"
		}
		builder = builder.OrderByClause(field)
	}

	return builder, nil
}

type Pagination struct {
	From *uint64
	Size *uint64
}

func (p *Pagination) Query(field string, builder squirrel.SelectBuilder) (squirrel.SelectBuilder, error) {
	if p.From != nil && *p.From > 0 {
		builder = builder.Offset(*p.From)
	}

	if p.Size == nil {
		size := uint64(20) //nolint:gomnd
		p.Size = &size
	}
	if p.Size != nil {
		builder = builder.Limit(*p.Size)
	}

	return builder, nil
}

type FindOptions struct {
	SelectForUpdate *bool

	Pagination *Pagination
	SortBy     SortBys
}

func (f *FindOptions) SetSelectForUpdate(v bool) *FindOptions {
	f.SelectForUpdate = &v
	return f
}

func (f *FindOptions) SetPagination(from uint64, size uint64) *FindOptions {
	f.Pagination = &Pagination{From: &from, Size: &size}
	return f
}

func (f *FindOptions) SetSortBy(sortBy ...SortBy) *FindOptions {
	f.SortBy = sortBy
	return f
}

func (f *FindOptions) Query(field string, builder squirrel.SelectBuilder) (squirrel.SelectBuilder, error) {
	if f.Pagination != nil {
		builder, _ = f.Pagination.Query("", builder)
	}

	builder, _ = f.SortBy.Query("", builder)
	return builder, nil
}

func Find() *FindOptions {
	return &FindOptions{}
}

func MergeFindOptions(opts ...*FindOptions) *FindOptions {
	fo := Find()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.SelectForUpdate != nil {
			fo.SelectForUpdate = opt.SelectForUpdate
		}
		if opt.Pagination != nil {
			fo.Pagination = opt.Pagination
		}
		if opt.SortBy != nil {
			fo.SortBy = opt.SortBy
		}
	}
	return fo
}

var ErrArgsMustBePointer = fmt.Errorf("Args object should be pointer")

// Build add condition to db
func Build(pargs interface{}, builder squirrel.SelectBuilder, except ...string) (squirrel.SelectBuilder, error) {
	if reflect.ValueOf(pargs).Kind() != reflect.Ptr {
		return builder, ErrArgsMustBePointer
	}

	argsVal := reflect.ValueOf(pargs).Elem()
	argsType := argsVal.Type()

	exceptMap := make(map[string]bool)
	for _, e := range except {
		exceptMap[e] = true
	}

	query := builder
	var err error
	for i := 0; i < argsVal.NumField(); i++ {
		// TODO: name 从 tag 中取，backoff 到 ToSnake
		fieldName := argsType.Field(i).Name
		name := strcase.ToSnake(fieldName)
		if _, ok := exceptMap[name]; ok {
			log.Tracef("field %s skip", name)
			continue
		}
		if _, ok := exceptMap[fieldName]; ok {
			log.Tracef("field %s skip", fieldName)
			continue
		}

		argsValFieldN := argsVal.Field(i)
		if argsValFieldN.IsNil() {
			continue
		}
		if argsValFieldN.CanInterface() {
			if fieldVal, ok := argsValFieldN.Interface().(Query); ok {
				query, err = fieldVal.Query(name, query)
				if err != nil {
					return builder, err
				}
			}
		}
	}

	return query, nil
}
