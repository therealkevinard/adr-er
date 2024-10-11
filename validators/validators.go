package validators

import (
	"fmt"
	"strconv"

	"github.com/therealkevinard/adr-er/globals"
)

// StrLenValidator returns a func that ensures a string's len is within range.
func StrLenValidator(fieldLabel string, min, max int) func(string) error {
	return func(s string) error {
		if len(s) < min {
			return globals.ValidationError(fieldLabel, fmt.Sprintf("must have at least %d characters", min))
		}

		if len(s) > max {
			return globals.ValidationError(fieldLabel, fmt.Sprintf("must have at most %d characters", min))
		}

		return nil
	}
}

// Uint32Between returns a validation func that ensures a string is numeric and the num value is within range.
func Uint32Between(fieldLabel string, min, max int) func(string) error {
	return func(str string) error {
		v, err := strconv.Atoi(str)
		if err != nil {
			return globals.ValidationError(fieldLabel, fmt.Sprintf("%s must be a number", str))
		}

		if v < min || v > max {
			return globals.ValidationError(fieldLabel, fmt.Sprintf("must be in range [%d, %d]", min, max))
		}

		return nil
	}
}
