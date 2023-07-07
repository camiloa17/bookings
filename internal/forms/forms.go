package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors{},
	}
}

// Required checks if required form fields are filled in
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	return strings.TrimSpace(f.Get(field)) != ""
}

// Valid check if the form is valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, minLength int) bool {
	fieldValue := f.Get(field)
	fieldValueLength := len(fieldValue)
	if fieldValueLength < minLength {
		f.Errors.Add(field, fmt.Sprintf("The min length is %d characters. Current length %d", minLength, fieldValueLength))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
		return false
	}
	return true
}
