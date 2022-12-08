package query

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

type SingleStringQuery struct {
	Op string `json:"op"`
	T  string `json:"t"`
}

type StringQuery []SingleStringQuery

func (q StringQuery) Query(field string, builder squirrel.SelectBuilder) (squirrel.SelectBuilder, error) {
	query := builder
	if len(q) == 0 {
		return query, nil
	}

	var cond []squirrel.Sqlizer

	for _, singleQuery := range q {
		if singleQuery.Op == "<>" {
			singleQuery.Op = "!="
		}
		switch singleQuery.Op {
		case "=":
			cond = append(cond, squirrel.Eq{field: singleQuery.T})
		case "!=":
			cond = append(cond, squirrel.NotEq{field: singleQuery.T})
		case "like":
			cond = append(cond, squirrel.Like{field: fmt.Sprint("%", singleQuery.T, "%")})
		case "not like":
			cond = append(cond, squirrel.NotLike{field: fmt.Sprint("%", singleQuery.T, "%")})
		default:
			return query, ErrInvalidFilterOperator
		}
	}

	var newCond squirrel.Sqlizer
	if len(cond) == 1 {
		newCond = cond[0]
	} else {
		newCond = squirrel.Or(cond)
	}

	query = query.Where(newCond)
	return query, nil
}
