package bitrix24server

func validateListTasksArgs(start *int) error {
	return validateOptionalNonNegativeInt("start", start)
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
