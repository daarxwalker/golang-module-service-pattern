package service

import (
	"example/app/repository"
	"example/core"
)

type Container struct {
	repository     repository.Container
	provideService core.ProvideService
	errorService   core.ErrorService
}

func NewContainer() func(ps core.ProvideService) Container {
	return func(ps core.ProvideService) Container {
		return Container{
			repository:     repository.NewContainer(ps.GetDB()),
			provideService: ps,
			errorService:   core.NewErrorService(ps),
		}
	}
}

func (c Container) ProvideService() core.ProvideService {
	return c.provideService
}

func (c Container) ErrorService() core.ErrorService {
	return c.errorService
}

func (c Container) CategoryService() CategoryService {
	return categoryService{
		provideService:     c.provideService,
		errorService:       c.errorService,
		categoryRepository: c.repository.GetCategoryRepository(),
	}
}
