package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	InputErrors map[string]string
}

func (vldtr *Validator) AddFieldError(key, value string) {
	if vldtr.InputErrors == nil {
		vldtr.InputErrors = make(map[string]string)
	}
	if _, contains := vldtr.InputErrors[key]; !contains {
		vldtr.InputErrors[key] = value
	}
}

func (vldtr *Validator) Valid() bool {
	return len(vldtr.InputErrors) == 0
}

func (vldtr *Validator) CheckField(valid bool, key, value string) {
	if !valid {
		vldtr.AddFieldError(key, value)
	}
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, lnth int) bool {
	return utf8.RuneCountInString(value) <= lnth
}
