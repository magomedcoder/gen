package service

import "testing"

func TestParseToolCallsJSON_empty(t *testing.T) {
	got, err := parseToolCallsJSON("")
	if err != nil || len(got) != 0 {
		t.Fatalf("parse empty: %v %#v", err, got)
	}
}

func TestParseToolCallsJSON_openAI(t *testing.T) {
	raw := `[{"id":"call_1","type":"function","function":{"name":"weather","arguments":"{\"city\":\"Paris\"}"}}]`
	got, err := parseToolCallsJSON(raw)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0].ID != "call_1" || got[0].Function.Name != "weather" {
		t.Fatalf("got %#v", got)
	}

	if got[0].Function.Arguments.M == nil || got[0].Function.Arguments.M["city"] != "Paris" {
		t.Fatalf("args map %#v", got[0].Function.Arguments.M)
	}
}
