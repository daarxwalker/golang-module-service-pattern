package service

import (
	"example/app/model/categoryModel"
	"example/app/module"
	"example/app/repository"
	"example/core"
	"example/core/form"
)

type CategoryService interface {
	GetAll(param core.Param) []categoryModel.Category
	GetFormSelect() []form.SelectOption
	GetOne(ID int) categoryModel.Category
	GetDefault() categoryModel.Category
	CreateOne(data categoryModel.Category) categoryModel.Category
	UpdateOne(data categoryModel.Category) categoryModel.Category
	Remove(data []int) bool
}

type categoryService struct {
	provideService     core.ProvideService
	errorService       core.ErrorService
	categoryRepository repository.CategoryRepository
}

func (s categoryService) GetAll(param core.Param) []categoryModel.Category {
	result, err := s.categoryRepository.GetAll(param)
	s.errorService.Check(err, errorMessage().getAll(module.Category))
	return result
}

func (s categoryService) GetFormSelect() []form.SelectOption {
	result, err := s.categoryRepository.GetFormSelect()
	s.errorService.Check(err, errorMessage().getForm(module.Category))
	return result
}

func (s categoryService) GetOne(ID int) categoryModel.Category {
	data, err := s.categoryRepository.GetOne(ID)
	s.errorService.Check(err, errorMessage().getOne(module.Category))
	return data
}

func (s categoryService) GetDefault() categoryModel.Category {
	data, err := s.categoryRepository.GetDefault()
	s.errorService.Check(err, errorMessage().getOne(module.Category))
	return data
}

func (s categoryService) CreateOne(data categoryModel.Category) categoryModel.Category {
	result, err := s.categoryRepository.CreateOne(data)
	s.errorService.Check(err, errorMessage().createOne(module.Category))
	return result
}

func (s categoryService) UpdateOne(data categoryModel.Category) categoryModel.Category {
	result, err := s.categoryRepository.UpdateOne(data)
	s.errorService.Check(err, errorMessage().updateOne(module.Category))
	return result
}

func (s categoryService) Remove(data []int) bool {
	result, err := s.categoryRepository.Remove(data)
	s.errorService.Check(err, errorMessage().remove(module.Category))
	return result
}
