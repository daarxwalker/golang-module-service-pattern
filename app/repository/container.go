package repository

import "github.com/go-pg/pg/v10"

type Container interface {
	GetCategoryRepository() CategoryRepository
}

type container struct {
	db *pg.DB
}

func NewContainer(db *pg.DB) Container {
	return container{
		db,
	}
}

func (c container) GetCategoryRepository() CategoryRepository {
	return categoryRepository{c.db}
}
