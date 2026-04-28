package bitrix24server

import (
	"fmt"
	"time"
)

type ExecutionDriftReport struct {
	TaskID              int     `json:"task_id"`
	TimeEstimateSeconds int     `json:"time_estimate_seconds"`
	TimeSpentSeconds    int     `json:"time_spent_seconds"`
	UtilizationRatio    float64 `json:"utilization_ratio"`
	OverrunSeconds      int     `json:"overrun_seconds"`
	CommentsCount       int     `json:"comments_count"`
	LastCommentAt       string  `json:"last_comment_at,omitempty"`
	SilenceHours        int     `json:"silence_hours"`
	DriftLevel          string  `json:"drift_level"`
	Summary             string  `json:"summary"`
}

func buildExecutionDriftReport(task TaskSnapshot, comments []CommentSnapshot, now time.Time) ExecutionDriftReport {
	estimate := task.TimeEstimate
	spent := task.TimeSpent
	ratio := 0.0
	if estimate > 0 {
		ratio = float64(spent) / float64(estimate)
	}

	overrun := 0
	if spent > estimate && estimate > 0 {
		overrun = spent - estimate
	}

	lastComment := ""
	silenceHours := 0
	var latest time.Time
	for _, c := range comments {
		t := parseBitrixTime(c.CreatedAt)
		if t.After(latest) {
			latest = t
			lastComment = c.CreatedAt
		}
	}

	if !latest.IsZero() {
		silenceHours = int(now.Sub(latest).Hours())
		if silenceHours < 0 {
			silenceHours = 0
		}
	}

	level := "low"
	switch {
	case (estimate > 0 && ratio >= 1.3) || silenceHours >= 96:
		level = "high"
	case (estimate > 0 && ratio >= 1.0) || silenceHours >= 48:
		level = "medium"
	}

	summary := "Дрифт исполнения низкий."
	if level == "medium" {
		summary = "Есть признаки дрифта исполнения: нужен контроль следующего шага."
	} else if level == "high" {
		summary = "Высокий дрифт исполнения: требуется перепланирование и снятие блокеров."
	}

	return ExecutionDriftReport{
		TaskID:              task.ID,
		TimeEstimateSeconds: estimate,
		TimeSpentSeconds:    spent,
		UtilizationRatio:    ratio,
		OverrunSeconds:      overrun,
		CommentsCount:       len(comments),
		LastCommentAt:       lastComment,
		SilenceHours:        silenceHours,
		DriftLevel:          level,
		Summary:             summary,
	}
}

func executionDriftActions(r ExecutionDriftReport) []string {
	actions := make([]string, 0, 3)
	if r.DriftLevel == "high" {
		actions = append(actions, "Сделать срочное перепланирование по задаче и согласовать новый контрольный срок.")
	}

	if r.OverrunSeconds > 0 {
		actions = append(actions, fmt.Sprintf("Уточнить оценку трудозатрат: перерасход %d сек.", r.OverrunSeconds))
	}

	if r.SilenceHours >= 48 {
		actions = append(actions, "Запросить статус-апдейт и следующий шаг в комментариях.")
	}

	if len(actions) == 0 {
		actions = append(actions, "Поддерживать текущий темп выполнения и регулярный апдейт статуса.")
	}

	return actions
}
