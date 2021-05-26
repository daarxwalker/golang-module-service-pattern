package repository

import (
	"fmt"
	"strings"

	"example/core"
)

type OrderBuilder interface {
	mergeParam(values []core.OrderParam) OrderBuilder
	order(values ...string) OrderBuilder
	direction(direction string) OrderBuilder
	orderPrefix(prefix string) OrderBuilder
	slice() []string
	string() string
}

type orderBuilder struct {
	ordersPrefix    string
	ordersDirection string
	orders          []string
}

func newOrderBuilder(values ...string) OrderBuilder {
	instance := &orderBuilder{}
	instance.order(values...)
	return instance
}

func (b *orderBuilder) orderPrefix(prefix string) OrderBuilder {
	b.ordersPrefix = prefix
	return b
}

func (b *orderBuilder) direction(direction string) OrderBuilder {
	b.ordersDirection = direction
	return b
}

func (b *orderBuilder) mergeParam(values []core.OrderParam) OrderBuilder {
	for _, item := range values {
		b.orders = append(b.orders, fmt.Sprintf("%s %s", item.Key, item.Direction))
	}
	return b
}

func (b *orderBuilder) order(values ...string) OrderBuilder {
	orders := make([]string, 0)
	isOrdersDirectionExist := len(b.ordersDirection) > 0

	for _, item := range values {
		if isOrdersDirectionExist {
			orders = append(orders, fmt.Sprintf("%s %s", item, b.ordersDirection))
		} else {
			orders = append(orders, fmt.Sprintf("%s %s", item, "ASC"))
		}
	}

	if len(b.ordersPrefix) > 0 {
		ordersLen := len(orders)
		if ordersLen > 1 {
			for _, order := range orders {
				orders = append(orders, fmt.Sprintf("%s.%s", b.ordersPrefix, order))
			}
		}
		if ordersLen == 1 {
			orders = append(orders, fmt.Sprintf("%s.%s", b.ordersPrefix, orders[0]))
		}
		if len(orders) > 0 {
			b.orders = append(b.orders, orders...)
		}
	} else {
		if len(values) > 0 {
			b.orders = append(b.orders, values...)
		}
	}

	b.ordersDirection = ""
	return b
}

func (b orderBuilder) slice() []string {
	return b.orders
}

func (b orderBuilder) string() string {
	return strings.Join(b.orders, ",")
}
