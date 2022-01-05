package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	validatorCheck func(reflect.Value, string) error
	typeValidators map[string]validatorCheck
)

var validators = map[reflect.Kind]typeValidators{
	reflect.Int: {
		"min": func(value reflect.Value, strMin string) error {
			if min, err := strconv.Atoi(strMin); err == nil {
				intValue := int(value.Int())

				if intValue >= min {
					return nil
				}

				return ErrLessLength(strconv.Itoa(intValue), strMin)
			}

			return ErrUnsupportedValidationParam("min", strMin)
		},
		"max": func(value reflect.Value, strMax string) error {
			if max, err := strconv.Atoi(strMax); err == nil {
				intValue := int(value.Int())

				if intValue <= max {
					return nil
				}

				return ErrMoreLength(strconv.Itoa(intValue), strMax)
			}

			return ErrUnsupportedValidationParam("max", strMax)
		},
		"in": func(value reflect.Value, strInList string) error {
			inList := strings.Split(strInList, ",")
			intValue := int(value.Int())

			for _, strParam := range inList {
				if in, err := strconv.Atoi(strParam); err == nil {
					if intValue == in {
						return nil
					}

					continue
				}

				return ErrUnsupportedValidationParam("in", strParam)
			}

			return ErrInRange(strconv.Itoa(intValue), strInList)
		},
	},
	reflect.String: {
		"len": func(value reflect.Value, strLen string) error {
			if length, err := strconv.Atoi(strLen); err == nil {
				if len(value.String()) == length {
					return nil
				}

				return ErrLength(value.String(), strLen)
			}

			return ErrUnsupportedValidationParam("len", strLen)
		},
		"regexp": func(value reflect.Value, pattern string) error {
			matched, err := regexp.MatchString(pattern, value.String())
			if err != nil {
				return ErrUnsupportedValidationParam("regexp", pattern)
			}

			if !matched {
				return ErrMatch(value.String(), pattern)
			}

			return nil
		},
		"in": func(value reflect.Value, strInList string) error {
			inList := strings.Split(strInList, ",")
			stringValue := value.String()

			for _, param := range inList {
				if stringValue == param {
					return nil
				}

				continue
			}

			return ErrInRange(stringValue, strInList)
		},
	},
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type ValidationParams map[string]string

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for _, error := range v {
		sb.WriteString(error.Field)
		sb.WriteString(": ")
		sb.WriteString(error.Err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func Validate(v interface{}) error {
	structValue := reflect.ValueOf(v)

	if structValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < structValue.NumField(); i++ {
		fieldType := structValue.Type().Field(i)
		fieldTag := fieldType.Tag

		validationParamsString, ok := fieldTag.Lookup("validate")

		if !ok {
			continue
		}

		validationParams, err := parseValidationParams(validationParamsString)
		if err != nil {
			return err
		}

		fieldValue := structValue.Field(i)

		fieldErrors, err := validateField(fieldValue, fieldType, validationParams)
		validationErrors = append(validationErrors, fieldErrors...)

		if err != nil {
			return err
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return ErrValidation(validationErrors.Error())
}

func validateField(
	fieldValue reflect.Value,
	fieldType reflect.StructField,
	validationParams ValidationParams,
) (ValidationErrors, error) {
	validationErrors := make(ValidationErrors, 0)

	if fieldValue.Kind() == reflect.Slice {
		for i := 0; i < fieldValue.Len(); i++ {
			if err := validateValue(fieldValue.Index(i), validationParams); err != nil {
				if errors.Is(err, ErrProgram) {
					return nil, err
				}

				validationErrors = append(validationErrors, ValidationError{fieldType.Name, err})
			}
		}

		return validationErrors, nil
	}

	if err := validateValue(fieldValue, validationParams); err != nil {
		if errors.Is(err, ErrProgram) {
			return nil, err
		}

		validationErrors = append(validationErrors, ValidationError{fieldType.Name, err})
	}

	return validationErrors, nil
}

func parseValidationParams(paramsString string) (ValidationParams, error) {
	params := strings.Split(paramsString, "|")

	paramsMap := make(ValidationParams)
	for _, param := range params {
		keyValue := strings.Split(param, ":")

		if len(keyValue) != 2 {
			return nil, ErrWrongValidator(param)
		}

		paramsMap[keyValue[0]] = keyValue[1]
	}

	return paramsMap, nil
}

func validateValue(value reflect.Value, validationParams ValidationParams) error {
	typeValidators, ok := validators[value.Kind()]

	if !ok {
		return ErrUnsupportedType(value.Kind())
	}

	// Сначала убеждаемся, что нет ошибок в написании валидаторов
	for validatorName := range validationParams {
		_, ok := typeValidators[validatorName]

		if !ok {
			return ErrUnsupportedValidator(validatorName, value.Kind())
		}
	}

	for validatorName, param := range validationParams {
		validator := typeValidators[validatorName]

		if err := validator(value, param); err != nil {
			return err
		}
	}

	return nil
}
