package validator

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

func CollectErrors(err error) []string {
	validationErrors := make([]string, 0)
	if err == nil {
		return validationErrors
	}
	var allErrs govalidator.Errors
	if errors.As(err, &allErrs) {
		for _, fld := range allErrs {
			validationErrors = append(validationErrors, fld.Error())
		}
	}
	return validationErrors
}
