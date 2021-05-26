package repository

type QueryBuilder interface {
	join(tableName string) JoinBuilder
	prefix(prefix string) ColumnBuilder
	column(values ...string) ColumnBuilder
	orderPrefix(prefix string) OrderBuilder
	order(values ...string) OrderBuilder
	direction(direction string) OrderBuilder
}

type queryBuilder struct {
}

func builder() QueryBuilder {
	return queryBuilder{}
}

func (b queryBuilder) join(tableName string) JoinBuilder {
	return newJoinBuilder(tableName)
}

func (b queryBuilder) prefix(prefix string) ColumnBuilder {
	instance := newColumnBuilder()
	instance.prefix(prefix)
	return instance
}

func (b queryBuilder) column(values ...string) ColumnBuilder {
	return newColumnBuilder(values...)
}

func (b queryBuilder) orderPrefix(prefix string) OrderBuilder {
	instance := newOrderBuilder()
	instance.orderPrefix(prefix)
	return instance
}

func (b queryBuilder) order(values ...string) OrderBuilder {
	return newOrderBuilder(values...)
}

func (b queryBuilder) direction(direction string) OrderBuilder {
	instance := newOrderBuilder()
	instance.direction(direction)
	return instance
}
