package form

type ValidationFunc func(map[string]interface{}) (frontend []error, backend error)

type Form interface {
	AddField(name string, field Field) error
	GetFields() []Field
	RegisterValidationFunctions(validationFuncs ...ValidationFunc)
	GetValidationFunctions() []ValidationFunc
	RegisterInitialValues(values map[string]interface{}) error
	Save(values map[string]HTMLType) (string, error)
}

func ValuesAreValid(form Form, values map[string]interface{}) ([]error, map[string][]error, error) {
	formErrs := make([]error, 0)
	fieldsErrs := make(map[string][]error)

	fields := form.GetFields()
	for _, field := range fields {
		fieldName := field.GetName()
		fieldValue, exists := values[fieldName]
		if !exists {
			fieldValue = nil
		}
		fieldErrs, err := FieldValueIsValid(field, fieldValue)
		if err != nil {
			return formErrs, fieldsErrs, err
		}
		fieldsErrs[fieldName] = fieldErrs
	}

	validationFuncs := form.GetValidationFunctions()
	for _, validationFunc := range validationFuncs {
		frontend, err := validationFunc(values)
		if err != nil {
			return formErrs, fieldsErrs, err
		}
		if frontend != nil {
			formErrs = append(formErrs, frontend...)
		}
	}

	return formErrs, fieldsErrs, nil
}

func GetCleanData(form Form, values map[string]HTMLType) (map[string]interface{}, error) {
	cleanValues := make(map[string]interface{})

	for _, field := range form.GetFields() {
		fieldName := field.GetName()
		fieldValue, exists := values[fieldName]
		if !exists {
			fieldValue = ""
		}
		cleanValue, err := field.HTMLTypeToGoType(fieldValue)
		if err != nil {
			return nil, err
		}
		cleanValues[fieldName] = cleanValue
	}

	return cleanValues, nil
}
