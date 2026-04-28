package streamtext

import "testing"

func TestStreamTextDelta_prefix(t *testing.T) {
	prev := "Привет"
	next := "Привет!"
	if got := StreamTextDelta(prev, next); got != "!" {
		t.Fatalf("got %q want %q", got, "!")
	}
}

func TestStreamTextDelta_emptyPrev(t *testing.T) {
	if got := StreamTextDelta("", "абв"); got != "абв" {
		t.Fatalf("got %q", got)
	}
}

func TestStreamTextDelta_shrink(t *testing.T) {
	prev := "Hello world"
	next := "Hello"
	if got := StreamTextDelta(prev, next); got != "" {
		t.Fatalf("shrink: got %q want empty", got)
	}
}

func TestStreamTextDelta_invalidUTF8Cleaned(t *testing.T) {
	got := StreamTextDelta("", "ab\xff\xfe")
	if got != "ab" {
		t.Fatalf("got %q want ab", got)
	}
}

func TestStreamTextDelta_runePrefixWhenBytesDiverge(t *testing.T) {
	prev := "x"
	next := "Привет"
	got := StreamTextDelta(prev, next)
	if got != "Привет" {
		t.Fatalf("got %q want full next", got)
	}
}
