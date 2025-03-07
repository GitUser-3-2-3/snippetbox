package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	InputErrors map[string]string
}

func (vdr *Validator) AddFieldError(key, value string) {
	if vdr.InputErrors == nil {
		vdr.InputErrors = make(map[string]string)
	}
	if _, contains := vdr.InputErrors[key]; !contains {
		vdr.InputErrors[key] = value
	}
}

func (vdr *Validator) Valid() bool {
	return len(vdr.InputErrors) == 0
}

func (vdr *Validator) CheckField(ok bool, key, value string) {
	if !ok {
		vdr.AddFieldError(key, value)
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
