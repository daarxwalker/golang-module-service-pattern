package repository

import (
	"example/core/form"
	"time"

	"github.com/go-pg/pg/v10"

	"example/app/model/categoryModel"
	"example/core"
)

type CategoryRepository interface {
	GetAll(p core.Param) ([]categoryModel.Category, error)
	GetFormSelect() ([]form.SelectOption, error)
	GetOne(ID int) (categoryModel.Category, error)
	GetDefault() (categoryModel.Category, error)
	CreateOne(data categoryModel.Category) (categoryModel.Category, error)
	UpdateOne(data categoryModel.Category) (categoryModel.Category, error)
	Remove(data []int) (bool, error)
}

type categoryRepository struct {
	db *pg.DB
}

func (r categoryRepository) baseColumns() []string {
	return []string{"name"}
}

func (r categoryRepository) columns() []string {
	return append([]string{"id", "created_at"}, r.baseColumns()...)
}

func (r categoryRepository) createColumns() []string {
	return append([]string{"vectors"}, r.baseColumns()...)
}

func (r categoryRepository) updateColumns() []string {
	return append([]string{"updated_at", "vectors"}, r.baseColumns()...)
}

func (r categoryRepository) GetAll(p core.Param) ([]categoryModel.Category, error) {
	var data []categoryModel.Category

	err := r.db.
		Model(&data).
		Column(r.columns()...).
		WhereGroup(fulltext(p.Fulltext, []string{alias.Category})).
		Order(order(p.Order)...).
		Limit(limit).
		Offset(p.Offset).
		Select()
	if check(err) {
		return data, err
	}

	return data, nil
}

func (r categoryRepository) GetFormSelect() ([]form.SelectOption, error) {
	var data []form.SelectOption

	err := r.db.
		Model().
		TableExpr(tableAs(table.Category, alias.Category)).
		ColumnExpr(
			builder().
				prefix(alias.Category).
				column(as("id", "value"), as("name", "label")).
				string(),
		).
		Order(builder().order("name").string()).
		Select(&data)
	if check(err) {
		return data, err
	}

	return data, nil
}

func (r categoryRepository) GetOne(ID int) (categoryModel.Category, error) {
	var data categoryModel.Category

	data.ID = ID

	err := r.db.
		Model(&data).
		Column(r.columns()...).
		WherePK().
		Limit(1).
		Select()
	if check(err) {
		return data, err
	}

	return data, nil
}

func (r categoryRepository) GetDefault() (categoryModel.Category, error) {
	var data categoryModel.Category

	err := r.db.
		Model(&data).
		Column(r.columns()...).
		Where(builder().column("name").placeholder(), "default").
		Limit(1).
		Select()
	if check(err) {
		return data, err
	}

	return data, nil
}

func (r categoryRepository) CreateOne(data categoryModel.Category) (categoryModel.Category, error) {
	data.Vectors = vectors(data.Name)

	_, err := r.db.
		Model(&data).
		Column(r.createColumns()...).
		Returning("id").
		Insert()
	if check(err) {
		return data, err
	}

	return data, nil
}

func (r categoryRepository) UpdateOne(data categoryModel.Category) (categoryModel.Category, error) {
	data.Vectors = vectors(data.Name)
	updateTime := time.Now()
	data.UpdatedAt = &updateTime

	_, err := r.db.
		Model(&data).
		Column(r.updateColumns()...).
		WherePK().
		Returning("id").
		UpdateNotZero()
	if check(err) {
		return data, err
	}

	return data, nil
}

func (r categoryRepository) Remove(data []int) (bool, error) {
	ok, err := remove(r.db, (*categoryModel.Category)(nil), data)
	if err != nil || !ok {
		return false, err
	}
	return true, nil
}
