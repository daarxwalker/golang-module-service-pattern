package form

import (
	"reflect"
)

type Form interface {
	Validate(validateData) FieldErrorMap
	SetField() fieldsProvider
	GetStructure() []field
	NotExist() FieldErrorMap
	SetFieldOptions(name string, options []SelectOption) Form
}

type form struct {
	fieldsProvider fieldsProvider
}

func New(data ...interface{}) Form {
	instance := &form{}
	if len(data) == 0 {
		instance.fieldsProvider = newFieldsProvider()
		return instance
	}

	value := reflect.ValueOf(data[0])
	if value.Len() == 0 {
		instance.fieldsProvider = newFieldsProvider()
		return instance
	}

	instance.fieldsProvider = newFieldsProvider(reflect.ValueOf(data[0]).Index(0).Interface())

	return instance
}

func (f *form) SetField() fieldsProvider {
	return f.fieldsProvider
}

func (f form) GetStructure() []field {
	return f.fieldsProvider.GetFields()
}

func (f *form) SetFieldOptions(name string, options []SelectOption) Form {
	f.fieldsProvider.updateFieldOptions(name, options)
	return f
}

func (f form) NotExist() FieldErrorMap {
	var err []fieldError
	return FieldErrorMap{"form": append(err, fieldError{Message: errorFormNotFound})}
}

func (f form) Validate(data validateData) FieldErrorMap {
	validator := newValidator(data)
	validator.validate(f.fieldsProvider.GetFields())
	return validator.getErrors()
}
