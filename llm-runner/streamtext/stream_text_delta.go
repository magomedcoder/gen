package streamtext

import "strings"

func StreamTextDelta(prev, next string) string {
	next = strings.ToValidUTF8(next, "")
	if next == "" {
		return ""
	}

	if prev == "" {
		return next
	}

	prev = strings.ToValidUTF8(prev, "")

	if len(next) >= len(prev) && strings.HasPrefix(next, prev) {
		return strings.ToValidUTF8(next[len(prev):], "")
	}

	pr := []rune(prev)
	nx := []rune(next)
	i := 0
	ml := len(pr)
	if len(nx) < ml {
		ml = len(nx)
	}

	for i < ml && pr[i] == nx[i] {
		i++
	}

	return strings.ToValidUTF8(string(nx[i:]), "")
}
