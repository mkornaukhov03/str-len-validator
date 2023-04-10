package strlenvalidator

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStrings := make([]string, len(v))
	for i, err := range v {
		errStrings[i] = err.Err.Error()
	}

	return strings.Join(errStrings, "\n")
}

func (v ValidationErrors) AddErr(err ValidationError) ValidationErrors {
	response := v
	if err.Err == nil {
		return response
	}
	if response == nil {
		response = make(ValidationErrors, 0)
	}
	response = append(response, err)
	return response
}

type validator interface {
	isValidStr(string) bool
}

type lenValidator struct {
	min int
	max int
}

func (l lenValidator) isValidStr(s string) bool {
	return l.min <= len(s) && len(s) < l.max
}

func getValidator(s string) (validator, error) {
	split := strings.Split(s, ":")
	if len(split) != 2 {
		return nil, ErrInvalidValidatorSyntax
	}
	cmd, args := split[0], split[1]
	argsSep := strings.Split(args, ",")
	if len(argsSep) != 2 {
		return nil, ErrInvalidValidatorSyntax
	}
	switch cmd {
	case "len":
		min, err := strconv.Atoi(argsSep[0])
		if err != nil {
			return nil, ErrInvalidValidatorSyntax
		}

		max, err := strconv.Atoi(argsSep[1])
		if err != nil {
			return nil, ErrInvalidValidatorSyntax
		}

		return lenValidator{min, max}, nil
	}
	return nil, ErrInvalidValidatorSyntax
}

func Validate(v any) error {
	inputValue := reflect.ValueOf(v)
	inputType := reflect.TypeOf(v)
	var errs ValidationErrors

	if inputValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < inputValue.NumField(); i++ {
		fieldValue := inputValue.Field(i)
		fieldType := inputType.Field(i)
		fieldTag := fieldType.Tag

		res := fieldTag.Get("validate")
		if res == "" {
			continue
		}
		if !fieldValue.CanInterface() {
			errs = errs.AddErr(ValidationError{ErrValidateForUnexportedFields})
			continue
		}
		validator, err := getValidator(res)
		if err != nil {
			errs = errs.AddErr(ValidationError{err})
			continue
		}
		var isValid bool
		if fieldType.Type == reflect.TypeOf("") {
			isValid = validator.isValidStr(fieldValue.Interface().(string))
		}
		if !isValid {
			err := fmt.Errorf("field %s: '%s' invalidation", fieldType.Name, res)
			errs = errs.AddErr(ValidationError{err})
		}
	}
	if errs == nil { // Because of interfaces internal implementation
		return nil
	}
	return errs
}
