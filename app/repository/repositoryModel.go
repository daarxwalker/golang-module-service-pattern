package repository

import (
	"example/app/model/categoryModel"
)

type repositoryModel struct {
	Category string
}

var table = repositoryModel{
	Category: categoryModel.CategoryTable,
}

var alias = repositoryModel{
	Category: categoryModel.CategoryAlias,
}
