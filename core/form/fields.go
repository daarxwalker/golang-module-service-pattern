package form

import (
	"encoding/json"
)

type fieldsProvider interface {
	Email(required ...bool) field
	Password(required ...bool) field
	Name(required ...bool) field
	Surname(required ...bool) field
	CountryCode(required ...bool) field
	City(required ...bool) field
	Street(required ...bool) field
	ZIP(required ...bool) field
	VAT(required ...bool) field
	Phone(required ...bool) field
	getField(name string) *field
	updateFieldOptions(name string, options []SelectOption)
	GetFields() []field
}

type fields struct {
	fields []field
	data   map[string]interface{}
}

func newFieldsProvider(data ...interface{}) fieldsProvider {
	instance := &fields{}
	instance.setData(data...)
	return instance
}

func getRequired(required []bool) bool {
	if len(required) == 0 {
		return false
	}
	return required[0]
}

func (f *fields) setData(data ...interface{}) {
	if len(data) > 0 {
		var defaultData map[string]interface{}
		bytes, err := json.Marshal(data[0])
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, &defaultData)
		if err != nil {
			panic(err)
		}
		f.data = defaultData
	}
}

func (f fields) getDefaultValue(name string, fieldType string) interface{} {
	if f.data[name] == nil {
		t := getFieldsTypes()
		switch fieldType {
		case t.Email:
			return ""
		case t.Password:
			return ""
		case t.Text:
			return ""
		case t.Number:
			return 0
		case t.Bool:
			return false
		case t.Select:
			return 0
		default:
			return ""
		}
	}
	return f.data[name]
}

func (f fields) getField(name string) *field {
	for i, field := range f.fields {
		if field.Name == name {
			return &f.fields[i]
		}
	}
	return &field{}
}

func (f *fields) updateFieldOptions(name string, options []SelectOption) {
	for i, field := range f.fields {
		if field.Name == name {
			f.fields[i].Options = options
		}
	}
}

func (f fields) GetFields() []field {
	return f.fields
}

func (f *fields) Email(required ...bool) field {
	name := Fields.Email
	newField := field{
		FieldType:    getFieldsTypes().Email,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Email),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         3,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Email,
				Value:         true,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         254,
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Password(required ...bool) field {
	name := Fields.Password
	newField := field{
		FieldType:    getFieldsTypes().Password,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Password),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         8,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         128,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Name(required ...bool) field {
	name := Fields.Name
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         2,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         255,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Surname(required ...bool) field {
	name := Fields.Surname
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         2,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         255,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) City(required ...bool) field {
	name := Fields.City
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         2,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         255,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Street(required ...bool) field {
	name := Fields.Street
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         2,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         255,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) CountryCode(required ...bool) field {
	name := Fields.CountryCode
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         3,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         3,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) ZIP(required ...bool) field {
	name := Fields.ZIP
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         1,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         30,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) VAT(required ...bool) field {
	name := Fields.VAT
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         1,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         30,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Phone(required ...bool) field {
	name := Fields.Phone
	newField := field{
		FieldType:    getFieldsTypes().Text,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Text),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().MinLength,
				Value:         1,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().MaxLength,
				Value:         30,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Category(required ...bool) field {
	name := Fields.Category
	newField := field{
		FieldType:    getFieldsTypes().Select,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Select),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().Min,
				Value:         1,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Distribution(required ...bool) field {
	name := Fields.Distribution
	newField := field{
		FieldType:    getFieldsTypes().Select,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Select),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().Min,
				Value:         1,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}

func (f *fields) Manufacturer(required ...bool) field {
	name := Fields.Manufacturer
	newField := field{
		FieldType:    getFieldsTypes().Select,
		Name:         name,
		DefaultValue: f.getDefaultValue(name, getFieldsTypes().Select),
		Validators: []fieldValidator{
			{
				ValidatorType: getFieldsValidatorsTypes().Min,
				Value:         1,
			},
			{
				ValidatorType: getFieldsValidatorsTypes().Required,
				Value:         getRequired(required),
			},
		},
	}

	f.fields = append(f.fields, newField)

	return newField
}
