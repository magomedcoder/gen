package template

import "testing"

func TestLoadAllPresets_embeddedPresets(t *testing.T) {
	presets, err := loadAllPresets()
	if err != nil {
		t.Fatal(err)
	}

	if len(presets) < 5 {
		t.Fatalf("expected several preset fingerprints, got %d", len(presets))
	}

	for _, p := range presets {
		if len(p.Bytes) == 0 {
			t.Fatalf("preset %q missing .gotmpl bytes", p.Name)
		}
	}
}
