package service

import (
	"slices"
	"testing"
)

func TestInferStopSequencesFromPrompt_chatML(t *testing.T) {
	p := "x" + chatMLImStart + "assistant\nhi" + chatMLImEnd
	got := inferStopSequencesFromPrompt(p)
	if !slices.Contains(got, chatMLImEnd) || !slices.Contains(got, chatMLImStart) {
		t.Fatalf("got %v", got)
	}
}

func TestInferStopSequencesFromPrompt_llama(t *testing.T) {
	p := "x" + llamaHdrStart + "assistant" + llamaHdrEnd + llamaEOT
	got := inferStopSequencesFromPrompt(p)
	if !slices.Contains(got, llamaEOT) {
		t.Fatalf("got %v", got)
	}
}

func TestInferStopSequencesFromPrompt_plain(t *testing.T) {
	if len(inferStopSequencesFromPrompt("User: hi\nAssistant:")) != 0 {
		t.Fatal("expected no inferred stops for plain prompt")
	}
}
