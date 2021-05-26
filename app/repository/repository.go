package repository

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"example/core/dd"

	"github.com/go-pg/pg/v10/types"

	"example/core"
	"example/core/helper/sliceHelper"
	"example/core/helper/stringHelper"
	"example/core/helper/vectorHelper"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

const (
	limit = 20
)

func check(err error) bool {
	return err != nil && err != pg.ErrNoRows
}

func as(column string, as string) string {
	return fmt.Sprintf("%s as %s", column, as)
}

func tableAs(table string, as string) string {
	return fmt.Sprintf("%s AS %s", table, as)
}

func prefix(prefix string, column string) string {
	return fmt.Sprintf("%s.%s", prefix, column)
}

func concat(values []string, dividers []string, as string) string {
	var result []string
	if len(values)-1 != len(dividers) {
		return ""
	}

	for i, item := range values {
		if len(values)-1 == i {
			result = append(result, item)
		} else {
			result = append(result, item)
			result = append(result, fmt.Sprintf("'%s'", dividers[i]))
		}
	}

	return fmt.Sprintf("concat(%s) as %s", strings.Join(result, ","), as)
}

func sum(column string) string {
	return fmt.Sprintf("sum(%s)", column)
}

func order(orders []core.OrderParam) []string {
	result := make([]string, len(orders))

	for i, item := range orders {
		result[i] = fmt.Sprintf("%s %s", stringHelper.SnakeCase(item.Key), strings.ToUpper(item.Direction))
	}

	return result
}

func calc(operand string, first string, second string) string {
	return fmt.Sprintf("%s %s %s", first, operand, second)
}

func in(value interface{}) types.ValueAppender {
	return pg.In(value)
}

func exist() string {
	return "exist(?)"
}

func returning(values ...string) string {
	return strings.Join(values, ",")
}

func subqueryPlaceholderAs(as string) string {
	return fmt.Sprintf("(?) AS %s", as)
}

func whereGroupIntSlice(slice []int, tableAlias string, column string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		for _, value := range slice {
			q = q.Where(builder().prefix(tableAlias).column(column).placeholder(), value)
		}
		return q, nil
	}
}

func optionalWhere(shouldCondition bool, condition string, params ...interface{}) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if shouldCondition {
			q = q.Where(condition, params...)
		}
		return q, nil
	}
}

func optionalIn(shouldWhere bool, placeholder string, values types.ValueAppender) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if shouldWhere {
			q = q.Where(placeholder, values)
		}
		return q, nil
	}
}

func fulltext(fulltext string, prefixes []string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if len(fulltext) > 0 {
			formattedFulltext := vectorHelper.Format(strings.ReplaceAll(fulltext, " ", "_"))
			prefixesLength := len(prefixes)
			for i, prefix := range prefixes {
				if i == 0 && prefixesLength > 1 {
					q = q.Where(fmt.Sprintf("%s.vectors @@ to_tsquery(?)", prefix), formattedFulltext+":*")
				} else {
					q = q.WhereOr(fmt.Sprintf("%s.vectors @@ to_tsquery(?)", prefix), formattedFulltext+":*")
				}
			}
		}
		return q, nil
	}
}

func remove(db *pg.DB, model interface{}, ids []int, column ...string) (bool, error) {
	c := "id"

	if len(column) > 0 {
		c = column[0]
	}

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	if _, err := tx.Model(model).
		Where(fmt.Sprintf("%s IN (?)", c), pg.In(ids)).
		Delete(); err != nil {
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func removeNoTx(db *pg.DB, model interface{}, ids []int, column ...string) (bool, error) {
	if _, err := db.Model(model).
		Where("id IN (?)", pg.In(ids)).
		Delete(); err != nil {
		return false, err
	}

	return true, nil
}

func vectors(data ...interface{}) types.ValueAppender {
	var r []string
	var vals []interface{}

	for _, item := range data {
		s := fmt.Sprintf("%v", item)
		if contains := sliceHelper.Contains(r, s); !contains && item != nil {
			for _, v := range regexp.MustCompile("[\\:\\,\\.\\_\\-\\s]+").Split(s, -1) {
				r = append(r, vectorHelper.Format(v))
			}
		}
	}

	return pg.SafeQuery("to_tsvector(?)", append(vals, strings.Join(r, " "))...)
}

func filtration(filters core.Filters, columns []string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		for key, filter := range filters {
			if !sliceHelper.Contains(columns, key) {
				continue
			}
			if filter.FilterType == core.FilterTypes.Select {
				q = q.Where(builder().column(key).placeholder(), filter.Value)
			}
			dd.Print(fmt.Sprintf("%T", filter.Value))
			if filter.FilterType == core.FilterTypes.Date && fmt.Sprintf("%T", filter.Value) == "[]time.Time" {
				value := filter.Value.([]time.Time)
				if len(value) == 1 {
					q = q.Where(builder().column(key).placeholderBetween(), value[0], time.Now())
				}
				if len(value) == 2 {
					q = q.Where(builder().column(key).placeholderBetween(), value[0], value[1])
				}
			}
		}
		return q, nil
	}
}
