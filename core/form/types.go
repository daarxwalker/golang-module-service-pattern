package form

func getFieldsTypes() fieldType {
	return fieldType{
		Email:    "email",
		Text:     "text",
		Password: "password",
		Number:   "number",
		Bool:     "bool",
		Select:   "select",
	}
}

func getFieldsValidatorsTypes() fieldValidatorType {
	return fieldValidatorType{
		Email:     "email",
		Min:       "min",
		Max:       "max",
		MinLength: "minLength",
		MaxLength: "maxLength",
		Required:  "required",
	}
}
