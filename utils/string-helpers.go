package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// PadValue pads the provided value with leading zeros up to width.
func PadValue(value, width int) string {
	return fmt.Sprintf("%0*d", width, value)
}

// SlugifySlice is a variadic wrapper for Slugify. provided strings are joined, and the result is slugified
func SlugifySlice(parts ...string) string {
	return Slugify(strings.Join(parts, "-"))
}

// Slugify converts a string into a normalized slug using several regex filters.
// The resulting slug is suitable for file name and other keying use-cases.
// TODO: with so many regexes, this func is a good candidate for fuzzing
func Slugify(input string) string {
	if input == "" {
		return input
	}

	slug := strings.ToLower(input) // tolower

	// reduce single or consecutive whitespace to a hyphen
	spaceRe := regexp.MustCompile(`\s+`)
	slug = spaceRe.ReplaceAllString(slug, "-")

	// keep only letters, numbers, and hyphens
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = re.ReplaceAllString(slug, "")

	// trim for tidiness
	slug = strings.Trim(slug, "-")

	return slug
}
