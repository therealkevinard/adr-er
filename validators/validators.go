package validators

import (
	"fmt"
	"strconv"
)

// StrLenValidator returns a func that ensures a string's len is within range
func StrLenValidator(fieldLabel string, min, max int) func(string) error {
	return func(s string) error {
		if len(s) < min {
			return fmt.Errorf("%s must have at least %d characters", fieldLabel, min)
		}
		if len(s) > max {
			return fmt.Errorf("%s must have at most %d characters", fieldLabel, max)
		}

		return nil
	}
}

// Uint32Between returns a validation func that ensures a string is numeric and the num value is within range
func Uint32Between(fieldLabel string, min, max int) func(string) error {
	return func(s string) error {
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("%s: %s is not numeric", fieldLabel, s)
		}
		if v < min || v > max {
			return fmt.Errorf("%s must be in range [%d, %d]", fieldLabel, min, max)
		}

		return nil
	}
}
