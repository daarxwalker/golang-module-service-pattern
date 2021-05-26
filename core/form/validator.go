package form

import (
	"fmt"
)

type validatorProvider interface {
	validate(fields []field)
	getErrors() FieldErrorMap
}

type validator struct {
	data   validateData
	errors FieldErrorMap
}

func newValidator(data validateData) validatorProvider {
	return &validator{
		data:   data,
		errors: make(FieldErrorMap),
	}
}

func (v validator) getErrors() FieldErrorMap {
	return v.errors
}

func (v *validator) checkErrorsNameExist(name string) {
	if v.errors[name] == nil {
		v.errors[name] = []fieldError{}
	}
}

func (v *validator) validateValidatorFormat(validatorValue interface{}, name string, validatorType string) bool {
	if fmt.Sprintf("%T", validatorValue) != validatorType {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorInvalidValidatorFormat,
		})
		return false
	}
	return true
}

func (v *validator) validateEmail(name string, value interface{}) {
	if fmt.Sprintf("%T", value) == "string" && !isEmailValid(value.(string)) {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorEmailInvalidFormat,
		})
	}
}

func (v *validator) validateMinLength(validatorValue interface{}, name string, value interface{}) {
	if !v.validateValidatorFormat(validatorValue, name, "int") {
		return
	}

	if fmt.Sprintf("%T", value) == "string" &&
		fmt.Sprintf("%T", validatorValue) == "int" &&
		len(value.(string)) < validatorValue.(int) {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorMinLength,
			Value:   validatorValue,
		})
	}
}

func (v *validator) validateMaxLength(validatorValue interface{}, name string, value interface{}) {
	if !v.validateValidatorFormat(validatorValue, name, "int") {
		return
	}

	if fmt.Sprintf("%T", value) == "string" &&
		fmt.Sprintf("%T", validatorValue) == "int" &&
		len(value.(string)) > validatorValue.(int) {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorMaxLength,
			Value:   validatorValue,
		})
	}
}

func (v *validator) validateMin(validatorValue interface{}, name string, value interface{}) {
	if !v.validateValidatorFormat(validatorValue, name, "int") {
		return
	}

	if fmt.Sprintf("%T", value) == "string" &&
		fmt.Sprintf("%T", validatorValue) == "int" &&
		value.(int) < validatorValue.(int) {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorMin,
			Value:   validatorValue,
		})
	}
}

func (v *validator) validateMax(validatorValue interface{}, name string, value interface{}) {
	if !v.validateValidatorFormat(validatorValue, name, "int") {
		return
	}

	if fmt.Sprintf("%T", value) == "string" &&
		fmt.Sprintf("%T", validatorValue) == "int" &&
		value.(int) > validatorValue.(int) {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorMax,
			Value:   validatorValue,
		})
	}
}

func (v *validator) validateRequired(validatorValue interface{}, name string, value interface{}) {
	if !v.validateValidatorFormat(validatorValue, name, "bool") {
		return
	}

	valueType := fmt.Sprintf("%T", value)
	if validatorValue.(bool) &&
		((valueType == "string" && len(value.(string)) == 0) ||
			(valueType == "int" && value.(int) == 0) ||
			(valueType == "bool" && value.(bool) == false)) {
		v.checkErrorsNameExist(name)
		v.errors[name] = append(v.errors[name], fieldError{
			Message: errorRequired,
			Value:   validatorValue,
		})
	}
}

func (v *validator) validate(fields []field) {
	vt := getFieldsValidatorsTypes()
	for name, value := range v.data {
		for _, item := range fields {
			if item.Name == name {
				for _, validator := range item.Validators {
					switch validator.ValidatorType {
					case vt.Email:
						v.validateEmail(name, value)
					case vt.MinLength:
						v.validateMinLength(validator.Value, name, value)
					case vt.MaxLength:
						v.validateMaxLength(validator.Value, name, value)
					case vt.Min:
						v.validateMin(validator.Value, name, value)
					case vt.Max:
						v.validateMax(validator.Value, name, value)
					case vt.Required:
						v.validateRequired(validator.Value, name, value)
					}
				}
			}
		}
	}

}
