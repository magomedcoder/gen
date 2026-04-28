package bitrix24server

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type BlockerSignal struct {
	CommentID    int    `json:"comment_id"`
	AuthorID     string `json:"author_id,omitempty"`
	DetectedAt   string `json:"detected_at,omitempty"`
	AgeHours     int    `json:"age_hours"`
	Message      string `json:"message"`
	Keyword      string `json:"keyword"`
	SeverityHint string `json:"severity_hint"`
}

func detectBlockerSignals(comments []CommentSnapshot, now time.Time) []BlockerSignal {
	keywords := []string{
		"блок", "заблок", "blocked", "blocker", "waiting", "жду", "ожида", "не могу", "проблем",
	}

	out := make([]BlockerSignal, 0, len(comments))
	for _, c := range comments {
		msgLower := strings.ToLower(strings.TrimSpace(c.Message))
		if msgLower == "" {
			continue
		}

		kw := ""
		for _, k := range keywords {
			if strings.Contains(msgLower, k) {
				kw = k
				break
			}
		}

		if kw == "" {
			continue
		}

		detectedAt := parseBitrixTime(c.CreatedAt)
		ageHours := 0
		if !detectedAt.IsZero() {
			ageHours = int(now.Sub(detectedAt).Hours())
			if ageHours < 0 {
				ageHours = 0
			}
		}

		severity := "medium"
		if ageHours >= 72 {
			severity = "high"
		} else if ageHours <= 24 {
			severity = "low"
		}

		out = append(out, BlockerSignal{
			CommentID:    c.ID,
			AuthorID:     c.AuthorID,
			DetectedAt:   c.CreatedAt,
			AgeHours:     ageHours,
			Message:      c.Message,
			Keyword:      kw,
			SeverityHint: severity,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].AgeHours != out[j].AgeHours {
			return out[i].AgeHours > out[j].AgeHours
		}

		return out[i].CommentID > out[j].CommentID
	})

	return out
}

func blockerOwners(signals []BlockerSignal) []string {
	set := map[string]struct{}{}
	for _, s := range signals {
		if strings.TrimSpace(s.AuthorID) == "" {
			continue
		}

		set[s.AuthorID] = struct{}{}
	}

	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}

	sort.Strings(out)

	return out
}

func blockerSummary(signals []BlockerSignal) string {
	if len(signals) == 0 {
		return "Блокеры в комментариях не обнаружены."
	}

	oldest := signals[0]
	return fmt.Sprintf("Обнаружено блокеров: %d. Самый старый: %dч назад (keyword=%s).", len(signals), oldest.AgeHours, oldest.Keyword)
}
