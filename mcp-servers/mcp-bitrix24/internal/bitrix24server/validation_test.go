package bitrix24server

import "testing"

func TestValidateOptionalNonNegativeInt(t *testing.T) {
	neg := -1
	if err := validateOptionalNonNegativeInt("start", &neg); err == nil {
		t.Fatalf("expected error for negative value")
	}

	zero := 0
	if err := validateOptionalNonNegativeInt("start", &zero); err != nil {
		t.Fatalf("unexpected error for zero: %v", err)
	}

	if err := validateOptionalNonNegativeInt("start", nil); err != nil {
		t.Fatalf("unexpected error for nil: %v", err)
	}
}

func TestValidateOptionalIntRange(t *testing.T) {
	inside := 10
	if err := validateOptionalIntRange("limit", &inside, 1, 50); err != nil {
		t.Fatalf("unexpected error for inside value: %v", err)
	}

	tooLow := 0
	if err := validateOptionalIntRange("limit", &tooLow, 1, 50); err == nil {
		t.Fatalf("expected error for low value")
	}

	tooHigh := 51
	if err := validateOptionalIntRange("limit", &tooHigh, 1, 50); err == nil {
		t.Fatalf("expected error for high value")
	}
}

func TestValidateOptionalEnum(t *testing.T) {
	if err := validateOptionalEnum("group_by", "responsible", "responsible", "creator", "status"); err != nil {
		t.Fatalf("unexpected error for allowed enum: %v", err)
	}

	if err := validateOptionalEnum("group_by", "STATUS", "responsible", "creator", "status"); err != nil {
		t.Fatalf("unexpected error for case-insensitive enum: %v", err)
	}

	if err := validateOptionalEnum("group_by", "", "responsible", "creator", "status"); err != nil {
		t.Fatalf("unexpected error for empty enum: %v", err)
	}

	if err := validateOptionalEnum("group_by", "team", "responsible", "creator", "status"); err == nil {
		t.Fatalf("expected error for unknown enum")
	}
}
