package categoryModule

import (
	"example/app/model/categoryModel"
	"example/app/service"
	"example/core"
)

type Controller interface {
	ProvideService() core.ProvideService
	CategoryService() service.CategoryService
}

func GetAll(c Controller) core.ControllerPayload {
	return c.CategoryService().GetAll(c.ProvideService().GetParams())
}

func GetOne(c Controller) core.ControllerPayload {
	return c.CategoryService().GetOne(c.ProvideService().GetParams().ID)
}

func CreateOne(c Controller) core.ControllerPayload {
	var data categoryModel.Category
	c.ProvideService().GetBody(&data)
	return c.CategoryService().CreateOne(data)
}

func UpdateOne(c Controller) core.ControllerPayload {
	var data categoryModel.Category
	c.ProvideService().GetBody(&data)
	return c.CategoryService().UpdateOne(data)
}

func Remove(c Controller) core.ControllerPayload {
	return c.CategoryService().Remove(c.ProvideService().GetParams().IDs)
}
