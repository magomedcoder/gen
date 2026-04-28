package bitrix24server

import "strings"

func taskID(task map[string]any) int {
	return numberLike(field(task, "id", "ID"))
}

func taskTitle(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "title", "TITLE"))
}

func taskStatusCode(task map[string]any) int {
	return numberLike(field(task, "status", "STATUS"))
}

func taskCreatedAtRaw(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "createdDate", "CREATED_DATE"))
}

func taskChangedAtRaw(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "changedDate", "CHANGED_DATE"))
}

func taskDeadlineRaw(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "deadline", "DEADLINE"))
}

func taskClosedAtRaw(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "closedDate", "CLOSED_DATE"))
}

func taskActivityAtRaw(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "activityDate", "ACTIVITY_DATE"))
}

func taskCreatedBy(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "createdBy", "CREATED_BY"))
}

func taskResponsibleID(task map[string]any) string {
	return strings.TrimSpace(stringField(task, "responsibleId", "RESPONSIBLE_ID"))
}

func taskPriority(task map[string]any) int {
	return numberLike(field(task, "priority", "PRIORITY"))
}

func taskTimeEstimate(task map[string]any) int {
	return numberLike(field(task, "timeEstimate", "TIME_ESTIMATE"))
}

func taskTimeSpent(task map[string]any) int {
	return numberLike(field(task, "timeSpentInLogs", "TIME_SPENT_IN_LOGS"))
}
