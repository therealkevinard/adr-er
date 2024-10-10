package globals

import "fmt"

// ValidationError returns either of  GenericInputValidationError or InputValidationError, depending on arguments.
func ValidationError(field, reason string) error {
	if field == "" && reason == "" {
		return GenericInputValidationError{}
	}

	return InputValidationError{
		Field:  field,
		Reason: reason,
	}
}

// GenericInputValidationError is used for generic invalid inputs.
type GenericInputValidationError struct{}

func (err GenericInputValidationError) Error() string { return "invalid input" }

// InputValidationError is used for precisely reporting invalid inputs
type InputValidationError struct {
	Field  string
	Reason string
}

func (err InputValidationError) Error() string {
	return fmt.Sprintf("%s failed validation: %s", err.Field, err.Reason)
}

// TemplateNotFoundError is used when no viable document template is available
type TemplateNotFoundError struct {
	Requested string
}

func (err TemplateNotFoundError) Error() string {
	return fmt.Sprintf("no template found for requested %s format", err.Requested)
}
