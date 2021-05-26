package repository

import "fmt"

type JoinBuilder interface {
	table(table string) JoinBuilder
	alias(alias string) JoinBuilder
	column(columns string) JoinBuilder
	on(prefix string, column string) string
}

type joinBuilder struct {
	tableName   string
	tableAlias  string
	tableColumn string
}

func newJoinBuilder(table string) JoinBuilder {
	return &joinBuilder{tableName: table}
}

func (b *joinBuilder) table(table string) JoinBuilder {
	b.tableName = table
	return b
}

func (b *joinBuilder) alias(alias string) JoinBuilder {
	b.tableAlias = alias
	return b
}

func (b *joinBuilder) column(columns string) JoinBuilder {
	b.tableColumn = columns
	return b
}

func (b joinBuilder) on(prefix string, column string) string {
	return fmt.Sprintf(
		"LEFT JOIN %s AS %s ON %s.%s = %s.%s",
		b.tableName,
		b.tableAlias,
		b.tableAlias,
		b.tableColumn,
		prefix,
		column,
	)
}
