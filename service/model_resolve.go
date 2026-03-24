package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func DisplayModelName(filename string) string {
	base := filepath.Base(filename)
	if len(base) > 5 && strings.EqualFold(base[len(base)-5:], ".gguf") {
		return base[:len(base)-5]
	}

	return base
}

func ListGGUFBasenames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if strings.EqualFold(filepath.Ext(name), ".gguf") {
			out = append(out, name)
		}
	}

	sort.Strings(out)

	return out, nil
}

func SortedDisplayModelNames(dir string) ([]string, error) {
	files, err := ListGGUFBasenames(dir)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(files))
	for i, f := range files {
		out[i] = DisplayModelName(f)
	}
	sort.Strings(out)

	return out, nil
}

func ResolveGGUFFile(modelsDir, userInput string) (canonical string, err error) {
	raw := strings.TrimSpace(userInput)
	if raw == "" {
		return "", fmt.Errorf("пустое имя модели")
	}

	base := filepath.Base(raw)
	if base == "." || base == string(filepath.Separator) {
		return "", fmt.Errorf("некорректное имя модели")
	}

	try := filepath.Join(modelsDir, base)
	if st, e := os.Stat(try); e == nil && !st.IsDir() {
		if strings.EqualFold(filepath.Ext(base), ".gguf") {
			return base, nil
		}
	}

	if filepath.Ext(base) == "" {
		cand := base + ".gguf"
		try2 := filepath.Join(modelsDir, cand)
		if _, e := os.Stat(try2); e == nil {
			return cand, nil
		}
	}

	files, err := ListGGUFBasenames(modelsDir)
	if err != nil {
		return "", err
	}

	for _, name := range files {
		if name == base {
			return name, nil
		}
	}

	for _, name := range files {
		if DisplayModelName(name) == base {
			return name, nil
		}
	}

	for _, name := range files {
		if strings.EqualFold(name, base) {
			return name, nil
		}
	}

	for _, name := range files {
		if strings.EqualFold(DisplayModelName(name), base) {
			return name, nil
		}
	}

	return "", fmt.Errorf("модель %q не найдена", userInput)
}
