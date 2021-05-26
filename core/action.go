package core

import (
	"example/core/form"
)

type Action interface {
	getClearFilesAfter() bool
	getController() Controller
	getForm() form.Form

	ClearFilesAfter() Action
	SetController(c Controller) Action
	ValidateForm(f form.Form) Action
}

type action struct {
	clearFilesAfter   bool
	form              form.Form
	formDefaultValues bool
	controller        Controller
}

func NewAction() Action {
	return &action{}
}

func (a action) getClearFilesAfter() bool {
	return a.clearFilesAfter
}

func (a action) getController() Controller {
	return a.controller
}

func (a action) getForm() form.Form {
	return a.form
}

func (a *action) ClearFilesAfter() Action {
	a.clearFilesAfter = true
	return a
}

func (a *action) SetController(c Controller) Action {
	a.controller = c
	return a
}

func (a *action) ValidateForm(f form.Form) Action {
	a.form = f
	return a
}
