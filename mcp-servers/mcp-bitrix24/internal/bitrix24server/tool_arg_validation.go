package bitrix24server

import (
	"fmt"
	"strings"
)

func validateListTasksArgs(start *int, filter map[string]any) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	return validateListFilterTaskIDConsistency(filter)
}

func validateListFilterTaskIDConsistency(filter map[string]any) error {
	if len(filter) == 0 {
		return nil
	}

	values := make(map[string]int)
	for _, key := range []string{"ID", "id", "=ID", "=id", "TASK_ID", "task_id", "taskId"} {
		raw, ok := filter[key]
		if !ok {
			continue
		}

		taskID, err := parseTaskID(raw)
		if err != nil || taskID <= 0 {
			continue
		}

		values[key] = taskID
	}

	if len(values) <= 1 {
		return nil
	}

	var expected int
	for _, v := range values {
		expected = v
		break
	}

	for key, v := range values {
		if v != expected {
			return fmt.Errorf("conflicting task id filters: %q=%d conflicts with other id fields", key, v)
		}
	}

	return nil
}

func validateAnalyzeTasksByQueryArgs(start, limit *int) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	return validateOptionalIntRange("limit", limit, 1, 50)
}

func validatePortfolioArgs(start, limit *int, groupBy string) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	if err := validateOptionalIntRange("limit", limit, 1, 50); err != nil {
		return err
	}

	return validateOptionalEnum("group_by", groupBy, "responsible", "creator", "status")
}

func validateExecutiveSummaryArgs(start, limit, periodDays *int) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	if err := validateOptionalIntRange("limit", limit, 1, 50); err != nil {
		return err
	}

	return validateOptionalIntRange("period_days", periodDays, 1, 30)
}

func validateSLAArgs(start, limit, soonHoursThreshold *int) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	if err := validateOptionalIntRange("limit", limit, 1, 50); err != nil {
		return err
	}

	return validateOptionalIntRange("soon_hours_threshold", soonHoursThreshold, 1, 168)
}

func validateWorkloadArgs(start, limit, overloadTasks *int) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	if err := validateOptionalIntRange("limit", limit, 1, 50); err != nil {
		return err
	}

	return validateOptionalIntRange("overload_tasks", overloadTasks, 1, 100)
}

func validateStatusTrendsArgs(start, limit, periodDays *int) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	if err := validateOptionalIntRange("limit", limit, 1, 50); err != nil {
		return err
	}

	return validateOptionalIntRange("period_days", periodDays, 1, 30)
}

func validateResponsiblePerformanceArgs(start, limit *int, responsibleID string) error {
	if err := validateOptionalNonNegativeInt("start", start); err != nil {
		return err
	}

	if err := validateOptionalIntRange("limit", limit, 1, 50); err != nil {
		return err
	}

	if strings.TrimSpace(responsibleID) == "" {
		return fmt.Errorf("responsible_id is required")
	}

	return nil
}
