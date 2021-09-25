package hw09structvalidator

import (
	"fmt"
	"reflect"
)

var (
	ErrStructValidation = fmt.Errorf("validation error")
	ErrProgram          = fmt.Errorf("program error")
	ErrNotStruct        = fmt.Errorf("struct should be provided")
)

// validation errors

func ErrLessLength(value string, min string) error {
	return fmt.Errorf("%q less than %q", value, min)
}

func ErrMoreLength(value string, max string) error {
	return fmt.Errorf("%q more than %q", value, max)
}

func ErrInRange(value string, in string) error {
	return fmt.Errorf("%q not in %q", value, in)
}

func ErrLength(value string, length string) error {
	return fmt.Errorf("%q length is not %q", value, length)
}

func ErrMatch(value string, pattern string) error {
	return fmt.Errorf("%q doesn't match %q", value, pattern)
}

func ErrValidation(text string) error {
	return fmt.Errorf("%w: %v", ErrStructValidation, text)
}

// program errors

func ErrUnsupportedType(kind reflect.Kind) error {
	return fmt.Errorf("%w: unsupported type %q", ErrProgram, kind)
}

func ErrUnsupportedValidator(name string, kind reflect.Kind) error {
	return fmt.Errorf("%w: unsupported validator %q for type %q", ErrProgram, name, kind)
}

func ErrWrongValidator(name string) error {
	return fmt.Errorf("%w: wrong validator %q", ErrProgram, name)
}

func ErrUnsupportedValidationParam(name string, value string) error {
	return fmt.Errorf("%w: unsupported %q value %q", ErrProgram, name, value)
}
