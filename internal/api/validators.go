package api

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var uuidSlice validator.Func = func(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < fl.Field().Len(); i++ {
		if !isValidUUID(fl.Field().Index(i).String()) {
			return false
		}
	}
	return true
}

func isValidUUID(s string) bool {
	parsed, err := uuid.Parse(s)
	return err == nil && parsed != uuid.Nil
}
