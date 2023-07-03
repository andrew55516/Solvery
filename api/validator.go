package api

import (
	"Solvery/util/tasks/task1"
	"github.com/go-playground/validator/v10"
)

var validTask1Input validator.Func = func(fl validator.FieldLevel) bool {
	if array, ok := fl.Field().Interface().([]int); ok {
		return task1.IsValidInput(array)
	}

	return false
}
