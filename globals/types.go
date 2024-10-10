package globals

// Validator implementations can validate themselves, returning errors.
// most/all of the types in the io pipeline should implement this interface.
type Validator interface {
	Validate() error
}
