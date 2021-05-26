package repository

import (
	"fmt"
	"strings"
)

type ColumnBuilder interface {
	column(values ...string) ColumnBuilder
	optionalColumn(shouldColumn bool, values ...string) ColumnBuilder
	prefix(prefix string) ColumnBuilder
	sumAs(as string) ColumnBuilder
	sum() ColumnBuilder
	placeholderBetween() string
	placeholder() string
	placeholderIn() string
	equal(column string) string
	filter(value bool) ColumnBuilder
	slice() []string
	string() string
	calc(operand string, as string) string
}

type columnBuilder struct {
	columnsPrefix string
	columns       []string
	columnSumAs   string
	columnSum     bool
	next          bool
	shouldFilter  bool
}

func newColumnBuilder(values ...string) ColumnBuilder {
	instance := &columnBuilder{}
	instance.column(values...)
	return instance
}

func (b columnBuilder) formatSingleColumn(column string) string {
	if b.columnSum {
		return fmt.Sprintf("sum(%s)", column)
	}
	if len(b.columnSumAs) > 0 {
		return fmt.Sprintf("sum(%s) as %s", column, b.columnSumAs)
	}
	return column
}

func (b *columnBuilder) sumAs(as string) ColumnBuilder {
	b.columnSumAs = as
	return b
}

func (b *columnBuilder) sum() ColumnBuilder {
	b.columnSum = true
	return b
}

func (b *columnBuilder) filter(next bool) ColumnBuilder {
	b.shouldFilter = true
	b.next = next
	return b
}

func (b *columnBuilder) resetSum() {
	b.columnSumAs = ""
	b.columnSum = false
}

func (b *columnBuilder) prefix(prefix string) ColumnBuilder {
	b.columnsPrefix = prefix
	return b
}

func (b *columnBuilder) afterColumn() {
	b.shouldFilter = false
	b.resetSum()
}

func (b *columnBuilder) optionalColumn(shouldColumn bool, values ...string) ColumnBuilder {
	if shouldColumn {
		b.column(values...)
	}
	return b
}

func (b *columnBuilder) column(values ...string) ColumnBuilder {
	if b.shouldFilter && !b.next {
		b.afterColumn()
		return b
	}

	if len(b.columnsPrefix) > 0 {
		var columns []string
		valuesLen := len(values)
		if valuesLen > 1 {
			for _, column := range values {
				columns = append(columns, fmt.Sprintf("%s.%s", b.columnsPrefix, column))
			}
		}
		if valuesLen == 1 {
			columns = append(columns, b.formatSingleColumn(fmt.Sprintf("%s.%s", b.columnsPrefix, values[0])))
		}
		if len(columns) > 0 {
			b.columns = append(b.columns, columns...)
		}
	} else {
		if len(values) > 0 {
			b.columns = append(b.columns, values...)
		}
	}

	b.afterColumn()

	return b
}

func (b *columnBuilder) equal(column string) string {
	if len(b.columnsPrefix) > 0 {
		return fmt.Sprintf("%s = %s.%s", b.columns[0], b.columnsPrefix, column)
	}
	return fmt.Sprintf("%s = %s", b.columns[0], column)
}

func (b *columnBuilder) placeholder() string {
	if len(b.columns) == 0 {
		return ""
	}

	return fmt.Sprintf("%s = ?", b.columns[0])
}

func (b *columnBuilder) placeholderIn() string {
	if len(b.columns) == 0 {
		return ""
	}

	return fmt.Sprintf("%s IN (?)", b.columns[0])
}

func (b *columnBuilder) placeholderBetween() string {
	if len(b.columns) == 0 {
		return ""
	}

	return fmt.Sprintf("%s ? between ?", b.columns[0])
}

func (b columnBuilder) slice() []string {
	return b.columns
}

func (b columnBuilder) calc(operand string, as string) string {
	if len(b.columns) < 2 {
		return strings.Join(b.columns, ",")
	}
	return fmt.Sprintf("(%s %s %s) AS %s", b.columns[0], operand, b.columns[1], as)
}

func (b columnBuilder) string() string {
	return strings.Join(b.columns, ",")
}
