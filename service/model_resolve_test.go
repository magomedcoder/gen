package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDisplayModelName(t *testing.T) {
	if g := DisplayModelName("Qwen-7B.Q4.gguf"); g != "Qwen-7B.Q4" {
		t.Fatalf("%q", g)
	}

	if g := DisplayModelName("LOWER.GGUF"); g != "LOWER" {
		t.Fatalf("%q", g)
	}
}

func TestResolveGGUFFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "MyModel-Q4.gguf"), []byte{0}, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveGGUFFile(dir, "MyModel-Q4")
	if err != nil || got != "MyModel-Q4.gguf" {
		t.Fatalf("got %q err %v", got, err)
	}

	got, err = ResolveGGUFFile(dir, "MyModel-Q4.gguf")
	if err != nil || got != "MyModel-Q4.gguf" {
		t.Fatalf("got %q err %v", got, err)
	}

	got, err = ResolveGGUFFile(dir, "mymodel-q4")
	if err != nil || got != "MyModel-Q4.gguf" {
		t.Fatalf("got %q err %v", got, err)
	}

	if _, err := ResolveGGUFFile(dir, "missing"); err == nil {
		t.Fatal("expected error")
	}
}

func TestSortedDisplayModelNames(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "b.gguf"), []byte{}, 0o644)
	_ = os.WriteFile(filepath.Join(dir, "a.gguf"), []byte{}, 0o644)
	got, err := SortedDisplayModelNames(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("%v", got)
	}
}
