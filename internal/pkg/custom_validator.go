package pkg

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasUppercase, _ := regexp.MatchString(`[A-Z]`, password)
	hasLowercase, _ := regexp.MatchString(`[a-z]`, password)
	hasNumber, _ := regexp.MatchString(`\d`, password)
	return hasUppercase && hasLowercase && hasNumber
}

func CategoryValidator(fl validator.FieldLevel) bool {
	category := fl.Field().String()
	const (
		CLIENT      = "client"
		NonBillable = "non-billable"
		SYSTEM      = "system"
	)
	validCategories := []string{CLIENT, NonBillable, SYSTEM}
	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}

func ProjectStartDateValidator(fl validator.FieldLevel) bool {
	startDate := fl.Field().Interface().(time.Time)
	return startDate.After(time.Now())
}
