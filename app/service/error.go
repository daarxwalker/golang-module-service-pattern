package service

import "fmt"

type message struct{}

func errorMessage() message {
	return message{}
}

func (m message) getAll(module string) string {
	return fmt.Sprintf("error.%s.getAll", module)
}

func (m message) getOne(module string) string {
	return fmt.Sprintf("error.%s.getOne", module)
}

func (m message) getForm(module string) string {
	return fmt.Sprintf("error.%s.getForm", module)
}

func (m message) createOne(module string) string {
	return fmt.Sprintf("error.%s.createOne", module)
}

func (m message) updateOne(module string) string {
	return fmt.Sprintf("error.%s.updateOne", module)
}

func (m message) remove(module string) string {
	return fmt.Sprintf("error.%s.remove", module)
}

func (m message) specific(module string, message string) string {
	return fmt.Sprintf("error.%s.%s", module, message)
}
