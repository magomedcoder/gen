package service

import (
	"strings"
	"testing"
)

const testImEnd = "<|im_end|>"

func TestTrimTrailingStops(t *testing.T) {
	got := trimTrailingStops("hello"+testImEnd, []string{testImEnd})
	if got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestStopStreamFilter_fullStopBufferedThenDropped(t *testing.T) {
	var b strings.Builder
	f := newStopStreamFilter([]string{testImEnd}, func(s string) { b.WriteString(s) })
	f.push(testImEnd)
	f.flush()

	if strings.Contains(b.String(), testImEnd) {
		t.Fatalf("stop leaked: %q", b.String())
	}
}

func TestStopStreamFilter_textThenStop(t *testing.T) {
	var b strings.Builder
	f := newStopStreamFilter([]string{testImEnd}, func(s string) { b.WriteString(s) })
	f.push("Привет!")
	f.push(testImEnd)
	f.flush()
	got := b.String()
	if strings.Contains(got, testImEnd) {
		t.Fatalf("got %q", got)
	}

	if !strings.Contains(got, "Привет") {
		t.Fatalf("lost text: %q", got)
	}
}

func TestStopStreamFilter_russianNotSplit(t *testing.T) {
	var b strings.Builder
	f := newStopStreamFilter([]string{testImEnd}, func(s string) { b.WriteString(s) })
	for _, r := range "Привет" {
		f.push(string(r))
	}

	f.push(testImEnd)
	f.flush()
	got := b.String()
	if strings.Contains(got, testImEnd) {
		t.Fatalf("got %q", got)
	}

	if got != "Привет" {
		t.Fatalf("want Привет, got %q", got)
	}
}
