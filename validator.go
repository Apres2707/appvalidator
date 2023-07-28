package appvalidator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func WithCustom(v *validator.Validate) error {
	err := v.RegisterValidation("max_without", maxWithout, false)
	if err != nil {
		return fmt.Errorf("register validation: %w", err)
	}

	return nil
}
