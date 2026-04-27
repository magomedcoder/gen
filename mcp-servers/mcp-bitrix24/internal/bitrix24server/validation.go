package bitrix24server

import (
	"fmt"
	"strings"
)

func validateOptionalNonNegativeInt(field string, value *int) error {
	if value == nil {
		return nil
	}

	if *value < 0 {
		return fmt.Errorf("%s must be >= 0", field)
	}

	return nil
}

func validateOptionalIntRange(field string, value *int, min, max int) error {
	if value == nil {
		return nil
	}

	if *value < min || *value > max {
		return fmt.Errorf("%s must be in range [%d..%d]", field, min, max)
	}

	return nil
}

func validateOptionalEnum(field, value string, allowed ...string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return nil
	}

	for _, item := range allowed {
		if value == strings.ToLower(strings.TrimSpace(item)) {
			return nil
		}
	}

	return fmt.Errorf("%s must be one of: %s", field, strings.Join(allowed, ", "))
}
