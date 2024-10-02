package form

type FieldValidationFunc func(value interface{}) (frontend []error, backend error)

type Field interface {
	HTML() (string, error)
	GoTypeToHTMLType(value interface{}) (HTMLType, error)
	HTMLTypeToGoType(value HTMLType) (interface{}, error)
	RegisterValidationFunctions(validationFuncs ...FieldValidationFunc)
	GetValidationFunctions() []FieldValidationFunc
	GetName() string
	RegisterName(name string) error
	RegisterInitialValue(value interface{})
}

func FieldValueIsValid(field Field, value interface{}) ([]error, error) {
	validationFuncs := field.GetValidationFunctions()
	var errs []error
	for _, validationFunc := range validationFuncs {
		frontend, err := validationFunc(value)
		if err != nil {
			return errs, err
		}
		if frontend != nil {
			errs = append(errs, frontend...)
		}
	}
	return errs, nil
}
