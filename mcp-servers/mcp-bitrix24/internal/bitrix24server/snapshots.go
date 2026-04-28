package bitrix24server

import (
	"fmt"
	"sort"
	"strings"
)

type TaskSnapshot struct {
	ID            int            `json:"id"`
	Title         string         `json:"title"`
	StatusCode    int            `json:"status_code"`
	StatusLabel   string         `json:"status_label"`
	CreatedAt     string         `json:"created_at,omitempty"`
	ChangedAt     string         `json:"changed_at,omitempty"`
	DeadlineAt    string         `json:"deadline_at,omitempty"`
	ClosedAt      string         `json:"closed_at,omitempty"`
	ActivityAt    string         `json:"activity_at,omitempty"`
	CreatedBy     string         `json:"created_by,omitempty"`
	ResponsibleID string         `json:"responsible_id,omitempty"`
	Priority      int            `json:"priority,omitempty"`
	TimeEstimate  int            `json:"time_estimate,omitempty"`
	TimeSpent     int            `json:"time_spent,omitempty"`
	Raw           map[string]any `json:"raw,omitempty"`
}

type CommentSnapshot struct {
	ID        int            `json:"id"`
	TaskID    int            `json:"task_id"`
	AuthorID  string         `json:"author_id,omitempty"`
	CreatedAt string         `json:"created_at,omitempty"`
	Message   string         `json:"message,omitempty"`
	Raw       map[string]any `json:"raw,omitempty"`
}

type TaskTimelineEvent struct {
	At      string `json:"at,omitempty"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
	Source  string `json:"source"`
}

func normalizeTaskSnapshot(task map[string]any) TaskSnapshot {
	statusCode := taskStatusCode(task)
	return TaskSnapshot{
		ID:            taskID(task),
		Title:         taskTitle(task),
		StatusCode:    statusCode,
		StatusLabel:   statusLabel(statusCode),
		CreatedAt:     taskCreatedAtRaw(task),
		ChangedAt:     taskChangedAtRaw(task),
		DeadlineAt:    taskDeadlineRaw(task),
		ClosedAt:      taskClosedAtRaw(task),
		ActivityAt:    taskActivityAtRaw(task),
		CreatedBy:     taskCreatedBy(task),
		ResponsibleID: taskResponsibleID(task),
		Priority:      taskPriority(task),
		TimeEstimate:  taskTimeEstimate(task),
		TimeSpent:     taskTimeSpent(task),
		Raw:           task,
	}
}

func normalizeCommentSnapshots(comments []map[string]any, fallbackTaskID int) []CommentSnapshot {
	out := make([]CommentSnapshot, 0, len(comments))
	for _, c := range comments {
		taskID := numberLike(field(c, "taskId", "task_id", "TASK_ID"))
		if taskID <= 0 {
			taskID = fallbackTaskID
		}

		out = append(out, CommentSnapshot{
			ID:        numberLike(field(c, "id", "ID")),
			TaskID:    taskID,
			AuthorID:  strings.TrimSpace(stringField(c, "authorId", "AUTHOR_ID")),
			CreatedAt: strings.TrimSpace(stringField(c, "post_date", "POST_DATE", "createdDate", "CREATED_DATE", "dateCreate", "DATE_CREATE")),
			Message:   strings.TrimSpace(stringField(c, "postMessage", "POST_MESSAGE", "message", "MESSAGE")),
			Raw:       c,
		})
	}

	return out
}

func buildTaskTimeline(task TaskSnapshot, comments []CommentSnapshot) []TaskTimelineEvent {
	events := make([]TaskTimelineEvent, 0, 8+len(comments))

	if task.CreatedAt != "" {
		events = append(events, TaskTimelineEvent{
			At:      task.CreatedAt,
			Type:    "task_created",
			Title:   "Задача создана",
			Details: fmt.Sprintf("Статус: %s", task.StatusLabel),
			Source:  "task",
		})
	}

	if task.DeadlineAt != "" {
		events = append(events, TaskTimelineEvent{
			At:      task.DeadlineAt,
			Type:    "task_deadline",
			Title:   "Дедлайн задачи",
			Details: task.Title,
			Source:  "task",
		})
	}

	if task.ChangedAt != "" {
		events = append(events, TaskTimelineEvent{
			At:     task.ChangedAt,
			Type:   "task_changed",
			Title:  "Последнее изменение задачи",
			Source: "task",
		})
	}

	if task.ActivityAt != "" {
		events = append(events, TaskTimelineEvent{
			At:     task.ActivityAt,
			Type:   "task_activity",
			Title:  "Последняя активность по задаче",
			Source: "task",
		})
	}

	if task.ClosedAt != "" {
		events = append(events, TaskTimelineEvent{
			At:     task.ClosedAt,
			Type:   "task_closed",
			Title:  "Задача закрыта",
			Source: "task",
		})
	}

	for _, c := range comments {
		details := c.Message
		if details == "" {
			details = "Комментарий без текста"
		}

		events = append(events, TaskTimelineEvent{
			At:      c.CreatedAt,
			Type:    "comment_added",
			Title:   "Добавлен комментарий",
			Details: details,
			Source:  "comment",
		})
	}

	sort.SliceStable(events, func(i, j int) bool {
		ti := parseBitrixTime(events[i].At)
		tj := parseBitrixTime(events[j].At)
		if ti.IsZero() && tj.IsZero() {
			return events[i].Type < events[j].Type
		}

		if ti.IsZero() {
			return false
		}

		if tj.IsZero() {
			return true
		}
		return ti.After(tj)
	})

	return events
}
