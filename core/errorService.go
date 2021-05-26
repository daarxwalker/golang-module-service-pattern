package core

import "strings"

type ErrorService interface {
	Check(err error, message ...string)
	New(message string)
}

type errorService struct {
	provideService ProvideService
}

func NewErrorService(provideService ProvideService) ErrorService {
	return errorService{
		provideService: provideService,
	}
}

func (s errorService) verify() {
	if s.provideService.GetTransactionExist() {
		s.provideService.RollbackTransaction()
	}
	if s.provideService.GetProtectExist() {
		s.provideService.RemoveCookies()
	}
}

func (s errorService) Check(err error, message ...string) {
	if err != nil {
		s.verify()
		panic(
			createError(
				err,
				strings.Join(message, ","),
			).getFormatted(),
		)
	}
}

func (s errorService) New(message string) {
	s.verify()
	panic(
		createNewError(
			message,
		).getFormatted(),
	)
}
