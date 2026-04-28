package bitrix24server

import "testing"

func TestParseBitrixTime_SupportsCommonFormats(t *testing.T) {
	cases := []string{
		"2026-04-28T15:30:00+03:00",
		"2026-04-28T15:30:00.000+03:00",
		"2026-04-28T15:30:00+0300",
		"2026-04-28T15:30:00.000+0300",
		"2026-04-28T15:30:00Z",
		"2026-04-28T15:30:00.000Z",
		"2026-04-28 15:30:00",
		"2026-04-28 15:30",
		"2026-04-28T15:30:00",
		"2026-04-28T15:30",
		"2026-04-28",
	}
	for _, input := range cases {
		tm := parseBitrixTime(input)
		if tm.IsZero() {
			t.Fatalf("expected parsed time for input %q", input)
		}
	}
}

func TestParseBitrixTime_InvalidReturnsZero(t *testing.T) {
	if tm := parseBitrixTime("not-a-time"); !tm.IsZero() {
		t.Fatalf("expected zero time for invalid input, got %v", tm)
	}
}
