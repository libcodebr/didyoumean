package verifier

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	ErrStructIsNil = errors.New("struct is nil")
)

// VerifyInterface is an interface that defines the methods for the Verify struct
type VerifyInterface interface {
	Struct(s interface{}) error
	Slice(s []interface{}) error
}

// Verifier is a global instance of the Verify struct that can be used to validate structs
var Verifier = NewVerify()

// NewVerify returns a new instance of the Verify struct
func NewVerify() VerifyInterface {
	return &Verify{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Verify is a struct that contains a validator instance
type Verify struct {
	validate *validator.Validate
}

// Struct verifies a struct using the validator and returns an error if the struct is not valid
func (v *Verify) Struct(s interface{}) error {
	if s == nil {
		return ErrStructIsNil
	}

	err := v.validate.Struct(s)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}

		var errs string
		for _, err := range err.(validator.ValidationErrors) {
			errs += v.errorMessage(err) + "\n"
		}

		return fmt.Errorf(strings.ToLower(errs))
	}

	return nil
}

// Slice verifies a slice of structs using the validator and returns an error if the slice is not valid
func (v *Verify) Slice(s []interface{}) error {
	if s == nil {
		return ErrStructIsNil
	}

	for _, i := range s {
		if err := v.Struct(i); err != nil {
			return err
		}
	}
	
	return nil
}

// errorMessage returns the error message for a field error
func (v *Verify) errorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())
	}
}
