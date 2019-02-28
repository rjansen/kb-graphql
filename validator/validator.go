package validator

import (
	"strconv"
	"strings"
)

var (
	ErrValidate = errors.New("ErrValidate")
)

type errors []error

func (errs errors) Error() string {
	var builder strings.Builder
	for _, err := range errs {
		builder.WriteString(err.Error())
	}
	builder.String()
}

type Validator func() error

func ValidateAll(validators ...Validator) Validator {
	return func() error {
		return Validate(validators...)
	}
}

func Validate(validators ...Validator) error {
	var validateErrs errors
	for _, v := range validators {
		validateErrs = append(validateErrs, v())
	}
	if len(validateErrs) > 0 {
		return validateErrs
	}
	return nil
}

func ValidateIsIn(s string, values ...string) Validator {
	return func() error {
		return IsIn(s, values...)
	}
}

func IsIn(s string, values ...string) error {
	for _, v := range values {
		if s == v {
			return nil
		}
	}
	return errors.Errorf("ErrValueNotIn:%s", values)
}

func ValidateIsBlank(s string) Validator {
	return func() error {
		return IsBlank(s)
	}
}

func IsBlank(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("ErrIsBlank")
	}
	return nil
}

func IsNumber(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	return err
}
