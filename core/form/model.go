package form

type field struct {
	FieldType    string           `json:"fieldType"`
	Name         string           `json:"name"`
	DefaultValue interface{}      `json:"defaultValue"`
	Options      []SelectOption   `json:"options,omitempty"`
	Validators   []fieldValidator `json:"validators"`
}

type validateData = map[string]interface{}

type fieldValidatorType struct {
	Email     string
	Min       string
	Max       string
	MinLength string
	MaxLength string
	Required  string
}

type fieldValidator struct {
	ValidatorType string      `json:"validatorType"`
	Value         interface{} `json:"value"`
}

type fieldType struct {
	Email    string
	Text     string
	Password string
	Number   string
	Bool     string
	Select   string
}

type fieldError struct {
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

type SelectOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

type FieldErrorMap = map[string][]fieldError
