package bitrix24server

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type taskAnalyticsItem struct {
	Task        map[string]any
	Comments    []map[string]any
	RiskScore   int
	HasBlockers bool
	IsOverdue   bool
	NoDeadline  bool
	NoComments  bool
	LastComment time.Time
	Deadline    time.Time
	Title       string
	TaskID      int
	StatusCode  int
	StatusLabel string
}

type statusCounters struct {
	Open       int
	InProgress int
	Done       int
	Deferred   int
}

func runAnalyticsQuery(ctx context.Context, client *bitrixClient, query string, taskID *int, filter map[string]any, order map[string]any, start *int, limit *int, includeComments *bool) (string, error) {
	now := time.Now()
	queryNorm := strings.ToLower(strings.TrimSpace(query))

	withComments := false
	if includeComments != nil {
		withComments = *includeComments
	}

	maxTasks := 20
	if limit != nil {
		if *limit < 1 {
			maxTasks = 1
		} else if *limit > 50 {
			maxTasks = 50
		} else {
			maxTasks = *limit
		}
	}

	items := make([]taskAnalyticsItem, 0, maxTasks)
	if taskID != nil && *taskID > 0 {
		task, err := loadTask(ctx, client, *taskID)
		if err != nil {
			return "", err
		}

		comments := []map[string]any(nil)
		if withComments {
			comments = loadTaskCommentsSoft(ctx, client, *taskID)
		}
		items = append(items, buildAnalyticsItem(task, comments, now))
	} else {
		tasks, err := loadTaskList(ctx, client, filter, order, start, maxTasks)
		if err != nil {
			return "", err
		}

		for _, task := range tasks {
			id := numberLike(field(task, "id", "ID"))
			if id <= 0 {
				continue
			}

			comments := []map[string]any(nil)
			if withComments {
				comments = loadTaskCommentsSoft(ctx, client, id)
			}
			items = append(items, buildAnalyticsItem(task, comments, now))
		}
	}

	if len(items) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	return renderAnalyticsAnswer(queryNorm, items, now), nil
}

func loadTask(ctx context.Context, client *bitrixClient, taskID int) (map[string]any, error) {
	resp, err := client.call(ctx, "tasks.task.get", map[string]any{
		"taskId": taskID,
		"select": []string{
			"ID", "TITLE", "STATUS", "CREATED_DATE", "CHANGED_DATE", "DEADLINE", "CREATED_BY", "RESPONSIBLE_ID", "PRIORITY", "TIME_ESTIMATE", "TIME_SPENT_IN_LOGS", "ACTIVITY_DATE",
		},
	})
	if err != nil {
		return nil, err
	}

	return extractTask(resp)
}

func loadTaskComments(ctx context.Context, client *bitrixClient, taskID int) ([]map[string]any, error) {
	resp, err := client.callTaskCommentItemGetList(ctx, taskID, map[string]any{"POST_DATE": "desc"}, nil)
	if err != nil {
		return nil, err
	}

	return extractComments(resp), nil
}

func loadTaskCommentsSoft(ctx context.Context, client *bitrixClient, taskID int) []map[string]any {
	comments, err := loadTaskComments(ctx, client, taskID)
	if err == nil {
		return comments
	}
	if isIgnorableCommentError(err) {
		log.Printf("[b24-mcp] analytics soft-skip comments task_id=%d err=%v", taskID, err)
		incSoftSkip()
		return nil
	}
	log.Printf("[b24-mcp] analytics comments unavailable task_id=%d err=%v", taskID, err)
	return nil
}

func isIgnorableCommentError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "task.commentitem.getlist") || strings.Contains(msg, "tasks_error_exception_#8") || strings.Contains(msg, "action_failed_to_be_processed") || strings.Contains(msg, "access denied")
}

func loadTaskList(ctx context.Context, client *bitrixClient, filter, order map[string]any, start *int, limit int) ([]map[string]any, error) {
	payload := map[string]any{
		"select": []string{
			"ID", "TITLE", "STATUS", "CREATED_DATE", "CHANGED_DATE", "DEADLINE", "CREATED_BY", "RESPONSIBLE_ID", "PRIORITY", "TIME_ESTIMATE", "TIME_SPENT_IN_LOGS", "ACTIVITY_DATE",
		},
	}

	if filter != nil {
		payload["filter"] = filter
	}

	if order != nil {
		payload["order"] = order
	}

	if start != nil {
		payload["start"] = *start
	}

	resp, err := client.call(ctx, "tasks.task.list", payload)
	if err != nil {
		return nil, err
	}

	result, _ := resp["result"].(map[string]any)
	rawTasks, _ := result["tasks"].([]any)
	all := toMapSlice(rawTasks)
	if len(all) > limit {
		all = all[:limit]
	}

	return all, nil
}

func buildAnalyticsItem(task map[string]any, comments []map[string]any, now time.Time) taskAnalyticsItem {
	id := numberLike(field(task, "id", "ID"))
	statusCode := numberLike(field(task, "status", "STATUS"))
	createdAt := parseBitrixTime(stringField(task, "createdDate", "CREATED_DATE"))
	deadline := parseBitrixTime(stringField(task, "deadline", "DEADLINE"))
	lastComment := lastCommentTime(comments)
	hasBlockers := scoreMentionsBlockers(comments)
	noComments := len(comments) == 0
	noDeadline := deadline.IsZero()
	isOverdue := !deadline.IsZero() && deadline.Before(now) && statusCode != 5 && statusCode != 7
	timeEstimate := numberLike(field(task, "timeEstimate", "TIME_ESTIMATE"))
	timeSpent := numberLike(field(task, "timeSpentInLogs", "TIME_SPENT_IN_LOGS"))

	score, _, _ := evaluateRisk(RiskInput{
		Now:           now,
		StatusCode:    statusCode,
		CreatedAt:     createdAt,
		Deadline:      deadline,
		LastComment:   lastComment,
		CommentsCount: len(comments),
		HasBlockers:   hasBlockers,
		TimeEstimate:  timeEstimate,
		TimeSpent:     timeSpent,
	}, defaultRiskScoringConfig())

	return taskAnalyticsItem{
		Task:        task,
		Comments:    comments,
		RiskScore:   score,
		HasBlockers: hasBlockers,
		IsOverdue:   isOverdue,
		NoDeadline:  noDeadline,
		NoComments:  noComments,
		LastComment: lastComment,
		Deadline:    deadline,
		Title:       stringField(task, "title", "TITLE"),
		TaskID:      id,
		StatusCode:  statusCode,
		StatusLabel: statusLabel(statusCode),
	}
}

func renderAnalyticsAnswer(query string, items []taskAnalyticsItem, now time.Time) string {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].RiskScore > items[j].RiskScore
	})

	var filtered []taskAnalyticsItem
	switch {
	case strings.Contains(query, "просроч"):
		for _, it := range items {
			if it.IsOverdue {
				filtered = append(filtered, it)
			}
		}
	case strings.Contains(query, "без дедлайн"):
		for _, it := range items {
			if it.NoDeadline {
				filtered = append(filtered, it)
			}
		}
	case strings.Contains(query, "блокер") || strings.Contains(query, "blocked") || strings.Contains(query, "риск"):
		for _, it := range items {
			if it.HasBlockers || it.RiskScore >= 4 {
				filtered = append(filtered, it)
			}
		}
	case strings.Contains(query, "без комментар") || strings.Contains(query, "тишин") || strings.Contains(query, "активност"):
		for _, it := range items {
			if it.NoComments || (!it.LastComment.IsZero() && now.Sub(it.LastComment) > 72*time.Hour) {
				filtered = append(filtered, it)
			}
		}
	default:
		filtered = items
	}
	if len(filtered) == 0 {
		return "По этому аналитическому запросу совпадений не найдено."
	}

	high := 0
	overdue := 0
	blockers := 0
	for _, it := range filtered {
		if it.RiskScore >= 5 {
			high++
		}

		if it.IsOverdue {
			overdue++
		}

		if it.HasBlockers {
			blockers++
		}
	}

	lines := []string{
		fmt.Sprintf("Запрос: %s", query),
		fmt.Sprintf("Найдено задач: %d", len(filtered)),
		fmt.Sprintf("Критичный риск: %d | Просрочено: %d | С блокерами: %d", high, overdue, blockers),
		"",
		"Топ задач по риску:",
	}

	top := filtered
	if len(top) > 10 {
		top = top[:10]
	}

	for _, it := range top {
		risk := "низкий"
		if it.RiskScore >= 5 {
			risk = "высокий"
		} else if it.RiskScore >= 2 {
			risk = "средний"
		}

		deadlineText := "без дедлайна"
		if !it.Deadline.IsZero() {
			deadlineText = it.Deadline.Format(time.RFC3339)
		}

		lines = append(lines, fmt.Sprintf("- #%d %s | статус: %s | риск: %s | дедлайн: %s | комментариев: %d", it.TaskID, emptyIfBlank(it.Title, "(без названия)"), it.StatusLabel, risk, deadlineText, len(it.Comments)))
	}

	lines = append(lines, "")
	lines = append(lines, "Рекомендованные действия:")
	if overdue > 0 {
		lines = append(lines, "- Срочно пересогласовать сроки и план закрытия по просроченным задачам.")
	}

	if blockers > 0 {
		lines = append(lines, "- Разобрать блокеры: назначить владельца каждого препятствия и дату снятия.")
	}

	lines = append(lines, "- Для задач без свежих комментариев запросить статус-апдейт и следующий контрольный шаг.")
	lines = append(lines, "")
	lines = append(lines, "=== Вывод ===")
	lines = append(lines, analyticsQueryConclusion(len(filtered), overdue, blockers, high))

	return strings.Join(lines, "\n")
}

func emptyIfBlank(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}

	return v
}

func analyticsQueryConclusion(total, overdue, blockers, highRisk int) string {
	if total == 0 {
		return formatConclusion("стабильно", "совпадений нет", "уточнить фильтр/запрос", "в плановом порядке")
	}

	if overdue > 0 || highRisk > 0 {
		return formatConclusion(
			"критично",
			fmt.Sprintf("просрочено %d, высокий риск %d", overdue, highRisk),
			"срочная стабилизация и перепланирование критичных задач",
			"сегодня",
		)
	}

	if blockers > 0 {
		return formatConclusion(
			"под контролем",
			fmt.Sprintf("блокеры: %d", blockers),
			"снять блокеры и закрепить владельцев",
			"в течение 1 рабочего дня",
		)
	}

	return formatConclusion("стабильно", "критичных отклонений нет", "поддерживать регулярный контроль", "по текущему регламенту")
}

func runPortfolioAnalytics(ctx context.Context, client *bitrixClient, filter map[string]any, order map[string]any, start *int, limit *int, includeComments *bool, groupBy string) (string, error) {
	now := time.Now()
	maxTasks := 30
	if limit != nil {
		if *limit < 1 {
			maxTasks = 1
		} else if *limit > 50 {
			maxTasks = 50
		} else {
			maxTasks = *limit
		}
	}

	withComments := false
	if includeComments != nil {
		withComments = *includeComments
	}

	tasks, err := loadTaskList(ctx, client, filter, order, start, maxTasks)
	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	items := make([]taskAnalyticsItem, 0, len(tasks))
	for _, task := range tasks {
		id := numberLike(field(task, "id", "ID"))
		if id <= 0 {
			continue
		}

		comments := []map[string]any(nil)
		if withComments {
			comments = loadTaskCommentsSoft(ctx, client, id)
		}

		items = append(items, buildAnalyticsItem(task, comments, now))
	}

	if len(items) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	normalizedGroupBy := strings.ToLower(strings.TrimSpace(groupBy))
	if normalizedGroupBy == "" {
		normalizedGroupBy = "responsible"
	}

	type groupStat struct {
		Name     string
		Total    int
		Overdue  int
		HighRisk int
		Blockers int
	}
	stats := map[string]*groupStat{}

	totalOverdue := 0
	totalHighRisk := 0
	totalBlockers := 0

	for _, it := range items {
		key := groupKey(it.Task, it.StatusLabel, normalizedGroupBy)
		if strings.TrimSpace(key) == "" {
			key = "(не указан)"
		}

		s, ok := stats[key]
		if !ok {
			s = &groupStat{Name: key}
			stats[key] = s
		}

		s.Total++
		if it.IsOverdue {
			s.Overdue++
			totalOverdue++
		}

		if it.RiskScore >= 5 {
			s.HighRisk++
			totalHighRisk++
		}

		if it.HasBlockers {
			s.Blockers++
			totalBlockers++
		}
	}

	rows := make([]groupStat, 0, len(stats))
	for _, s := range stats {
		rows = append(rows, *s)
	}

	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Overdue != rows[j].Overdue {
			return rows[i].Overdue > rows[j].Overdue
		}

		if rows[i].HighRisk != rows[j].HighRisk {
			return rows[i].HighRisk > rows[j].HighRisk
		}

		return rows[i].Total > rows[j].Total
	})

	sort.SliceStable(items, func(i, j int) bool { return items[i].RiskScore > items[j].RiskScore })

	lines := []string{
		fmt.Sprintf("Портфель задач (group_by=%s)", normalizedGroupBy),
		fmt.Sprintf("Всего задач: %d", len(items)),
		fmt.Sprintf("Просрочено: %d | Высокий риск: %d | С блокерами: %d", totalOverdue, totalHighRisk, totalBlockers),
		"",
		"Сводка по группам:",
	}

	for _, row := range rows {
		lines = append(lines, fmt.Sprintf("- %s: всего=%d, просрочено=%d, высокий риск=%d, блокеры=%d",
			row.Name, row.Total, row.Overdue, row.HighRisk, row.Blockers))
	}

	lines = append(lines, "")
	lines = append(lines, "Топ рискованных задач:")
	top := items
	if len(top) > 10 {
		top = top[:10]
	}

	for _, it := range top {
		risk := "низкий"
		if it.RiskScore >= 5 {
			risk = "высокий"
		} else if it.RiskScore >= 2 {
			risk = "средний"
		}

		lines = append(lines, fmt.Sprintf("- #%d %s | %s | риск=%s | просрочено=%t | блокеры=%t", it.TaskID, emptyIfBlank(it.Title, "(без названия)"), it.StatusLabel, risk, it.IsOverdue, it.HasBlockers))
	}

	lines = append(lines, "")
	lines = append(lines, "Рекомендации по портфелю:")
	lines = append(lines, "- Сначала разберите группы с максимальным числом просрочек.")
	lines = append(lines, "- Для задач высокого риска назначьте конкретные даты контрольных апдейтов.")
	lines = append(lines, "- Для задач с блокерами фиксируйте владельца снятия блокера и срок.")
	lines = append(lines, "")
	lines = append(lines, "=== Вывод ===")
	lines = append(lines, portfolioConclusion(len(items), totalOverdue, totalHighRisk, totalBlockers))

	return strings.Join(lines, "\n"), nil
}

func portfolioConclusion(total, overdue, highRisk, blockers int) string {
	if total == 0 {
		return formatConclusion("стабильно", "портфель пуст", "верифицировать фильтр выборки", "в плановом порядке")
	}

	if overdue > 0 || highRisk > 0 {
		return formatConclusion(
			"критично",
			fmt.Sprintf("просрочено %d, высокий риск %d", overdue, highRisk),
			"управленческое вмешательство по проблемным группам",
			"сегодня",
		)
	}

	if blockers > 0 {
		return formatConclusion("под контролем", fmt.Sprintf("блокеры: %d", blockers), "устранение блокеров по владельцам", "в течение 1 рабочего дня")
	}

	return formatConclusion("стабильно", "существенных рисков не выявлено", "поддерживать текущий контур управления", "по текущему регламенту")
}

func groupKey(task map[string]any, statusLabelValue string, groupBy string) string {
	switch groupBy {
	case "creator":
		return stringField(task, "createdBy", "CREATED_BY")
	case "status":
		return statusLabelValue
	default:
		return stringField(task, "responsibleId", "RESPONSIBLE_ID")
	}
}

func runExecutiveSummary(ctx context.Context, client *bitrixClient, filter map[string]any, order map[string]any, start *int, limit *int, periodDays *int, includeComments *bool) (string, error) {
	now := time.Now()
	maxTasks := 40
	if limit != nil {
		if *limit < 1 {
			maxTasks = 1
		} else if *limit > 50 {
			maxTasks = 50
		} else {
			maxTasks = *limit
		}
	}

	days := 7
	if periodDays != nil {
		if *periodDays < 1 {
			days = 1
		} else if *periodDays > 30 {
			days = 30
		} else {
			days = *periodDays
		}
	}
	currentStart := now.Add(-time.Duration(days) * 24 * time.Hour)
	previousStart := currentStart.Add(-time.Duration(days) * 24 * time.Hour)

	withComments := false
	if includeComments != nil {
		withComments = *includeComments
	}

	tasks, err := loadTaskList(ctx, client, filter, order, start, maxTasks)
	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	items := make([]taskAnalyticsItem, 0, len(tasks))
	for _, task := range tasks {
		id := numberLike(field(task, "id", "ID"))
		if id <= 0 {
			continue
		}

		comments := []map[string]any(nil)
		if withComments {
			comments = loadTaskCommentsSoft(ctx, client, id)
		}

		items = append(items, buildAnalyticsItem(task, comments, now))
	}

	if len(items) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	total := len(items)
	active := 0
	overdue := 0
	highRisk := 0
	noDeadline := 0
	blockers := 0
	changedCurrentItems := make([]taskAnalyticsItem, 0, len(items))

	currentChanged := 0
	previousChanged := 0
	currentCreated := 0
	previousCreated := 0

	for _, it := range items {
		if it.StatusCode != 5 && it.StatusCode != 7 {
			active++
		}

		if it.IsOverdue {
			overdue++
		}

		if it.RiskScore >= 5 {
			highRisk++
		}

		if it.NoDeadline {
			noDeadline++
		}

		if it.HasBlockers {
			blockers++
		}

		changedAt := parseBitrixTime(stringField(it.Task, "changedDate", "CHANGED_DATE"))
		createdAt := parseBitrixTime(stringField(it.Task, "createdDate", "CREATED_DATE"))
		if !changedAt.IsZero() {
			if changedAt.After(currentStart) {
				currentChanged++
				changedCurrentItems = append(changedCurrentItems, it)
			} else if changedAt.After(previousStart) && changedAt.Before(currentStart) {
				previousChanged++
			}
		}

		if !createdAt.IsZero() {
			if createdAt.After(currentStart) {
				currentCreated++
			} else if createdAt.After(previousStart) && createdAt.Before(currentStart) {
				previousCreated++
			}
		}
	}

	sort.SliceStable(items, func(i, j int) bool { return items[i].RiskScore > items[j].RiskScore })

	lines := []string{
		fmt.Sprintf("Executive summary за %d дн.", days),
		fmt.Sprintf("Охват: %d задач", total),
		fmt.Sprintf("Активные: %d | Просрочено: %d | Высокий риск: %d | Без дедлайна: %d | С блокерами: %d",
			active, overdue, highRisk, noDeadline, blockers),
		"",
		"Тренды vs предыдущий период:",
		fmt.Sprintf("- Измененных задач: %d (%s)", currentChanged, trendDelta(currentChanged, previousChanged)),
		fmt.Sprintf("- Новых задач: %d (%s)", currentCreated, trendDelta(currentCreated, previousCreated)),
		"",
		"What changed (топ изменений за текущий период):",
	}

	if len(changedCurrentItems) == 0 {
		lines = append(lines, "- Изменений за текущий период не зафиксировано.")
	} else {
		sort.SliceStable(changedCurrentItems, func(i, j int) bool {
			ai := parseBitrixTime(stringField(changedCurrentItems[i].Task, "changedDate", "CHANGED_DATE"))
			aj := parseBitrixTime(stringField(changedCurrentItems[j].Task, "changedDate", "CHANGED_DATE"))
			return ai.After(aj)
		})

		topChanged := changedCurrentItems
		if len(topChanged) > 5 {
			topChanged = topChanged[:5]
		}

		for _, it := range topChanged {
			changedAt := parseBitrixTime(stringField(it.Task, "changedDate", "CHANGED_DATE"))
			changedText := "n/a"
			if !changedAt.IsZero() {
				changedText = changedAt.Format(time.RFC3339)
			}

			lines = append(lines, fmt.Sprintf("- #%d %s | статус=%s | changed=%s | риск=%d",
				it.TaskID, emptyIfBlank(it.Title, "(без названия)"), it.StatusLabel, changedText, it.RiskScore))
		}
	}

	lines = append(lines,
		"",
		"Фокус руководителя (топ-10 рисков):",
	)

	top := items
	if len(top) > 10 {
		top = top[:10]
	}

	for _, it := range top {
		risk := "низкий"
		if it.RiskScore >= 5 {
			risk = "высокий"
		} else if it.RiskScore >= 2 {
			risk = "средний"
		}

		lines = append(lines, fmt.Sprintf("- #%d %s | %s | риск=%s | просрочено=%t | дедлайн=%s", it.TaskID, emptyIfBlank(it.Title, "(без названия)"), it.StatusLabel, risk, it.IsOverdue, execDeadlineText(it.Deadline)))
	}

	lines = append(lines, "")
	lines = append(lines, "Рекомендации:")
	lines = append(lines, "- Приоритизируйте просроченные задачи с высоким риском в ежедневном контроле.")
	lines = append(lines, "- Закройте пробелы в дедлайнах у активных задач.")
	lines = append(lines, "- По задачам с блокерами зафиксируйте владельца и срок снятия блока.")
	lines = append(lines, "")
	lines = append(lines, "=== Вывод ===")
	lines = append(lines, executiveConclusion(total, overdue, highRisk, noDeadline, currentChanged, previousChanged))

	return strings.Join(lines, "\n"), nil
}

func executiveConclusion(total, overdue, highRisk, noDeadline, currentChanged, previousChanged int) string {
	if total == 0 {
		return formatConclusion("стабильно", "недостаточно данных", "расширить охват выборки", "в плановом порядке")
	}

	trend := trendDelta(currentChanged, previousChanged)
	if overdue > 0 || highRisk > 0 {
		return formatConclusion(
			"критично",
			fmt.Sprintf("просрочено %d, высокий риск %d, динамика %s", overdue, highRisk, trend),
			"усилить контроль и пересобрать приоритеты портфеля",
			"сегодня",
		)
	}

	if noDeadline > 0 {
		return formatConclusion("под контролем", fmt.Sprintf("задачи без дедлайна: %d", noDeadline), "назначить сроки активным задачам", "в течение 1 рабочего дня")
	}

	return formatConclusion("стабильно", fmt.Sprintf("критичных рисков нет, динамика %s", trend), "поддерживать текущий управленческий ритм", "по текущему регламенту")
}

func trendDelta(current, previous int) string {
	delta := current - previous
	if delta > 0 {
		return fmt.Sprintf("+%d", delta)
	}

	if delta < 0 {
		return fmt.Sprintf("%d", delta)
	}

	return "0"
}

func execDeadlineText(deadline time.Time) string {
	if deadline.IsZero() {
		return "нет"
	}

	return deadline.Format(time.RFC3339)
}

func runSLASummary(ctx context.Context, client *bitrixClient, filter map[string]any, order map[string]any, start *int, limit *int, soonHoursThreshold *int, includeComments *bool) (string, error) {
	now := time.Now()
	maxTasks := 40
	if limit != nil {
		if *limit < 1 {
			maxTasks = 1
		} else if *limit > 50 {
			maxTasks = 50
		} else {
			maxTasks = *limit
		}
	}

	soonHours := 24
	if soonHoursThreshold != nil {
		if *soonHoursThreshold < 1 {
			soonHours = 1
		} else if *soonHoursThreshold > 168 {
			soonHours = 168
		} else {
			soonHours = *soonHoursThreshold
		}
	}

	withComments := false
	if includeComments != nil {
		withComments = *includeComments
	}

	tasks, err := loadTaskList(ctx, client, filter, order, start, maxTasks)
	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	items := make([]taskAnalyticsItem, 0, len(tasks))
	for _, task := range tasks {
		id := numberLike(field(task, "id", "ID"))
		if id <= 0 {
			continue
		}

		comments := []map[string]any(nil)
		if withComments {
			comments = loadTaskCommentsSoft(ctx, client, id)
		}

		items = append(items, buildAnalyticsItem(task, comments, now))
	}

	if len(items) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	type bucketed struct {
		item      taskAnalyticsItem
		urgency   int
		urgencyTx string
	}

	bucket := make([]bucketed, 0, len(items))
	noDeadline := 0
	completed := 0
	overdue := 0
	today := 0
	soon := 0

	for _, it := range items {
		if it.StatusCode == 5 || it.StatusCode == 7 {
			completed++
			continue
		}

		if it.Deadline.IsZero() {
			noDeadline++
			bucket = append(bucket, bucketed{
				item:      it,
				urgency:   3,
				urgencyTx: "P2: нет дедлайна",
			})
			continue
		}

		hoursLeft := it.Deadline.Sub(now).Hours()
		switch {
		case hoursLeft < 0:
			overdue++
			bucket = append(bucket, bucketed{item: it, urgency: 1, urgencyTx: fmt.Sprintf("P0: просрочено на %d ч.", int(-hoursLeft))})
		case hoursLeft <= 24:
			today++
			bucket = append(bucket, bucketed{item: it, urgency: 2, urgencyTx: fmt.Sprintf("P1: дедлайн <=24ч (осталось %d ч.)", int(hoursLeft))})
		case hoursLeft <= float64(soonHours):
			soon++
			bucket = append(bucket, bucketed{item: it, urgency: 3, urgencyTx: fmt.Sprintf("P2: скоро дедлайн (осталось %d ч.)", int(hoursLeft))})
		}
	}

	sort.SliceStable(bucket, func(i, j int) bool {
		if bucket[i].urgency != bucket[j].urgency {
			return bucket[i].urgency < bucket[j].urgency
		}

		return bucket[i].item.RiskScore > bucket[j].item.RiskScore
	})

	lines := []string{
		"SLA summary по задачам",
		fmt.Sprintf("Охват: %d задач | Завершено/отклонено: %d | Без дедлайна: %d", len(items), completed, noDeadline),
		fmt.Sprintf("Нарушение SLA (просрочено): %d | Критично сегодня: %d | Скоро дедлайн (%dч): %d", overdue, today, soonHours, soon),
		"",
		"Очередь реакции:",
	}

	top := bucket
	if len(top) > 15 {
		top = top[:15]
	}

	for _, b := range top {
		lines = append(lines, fmt.Sprintf("- %s | #%d %s | статус=%s | риск=%d | блокеры=%t",
			b.urgencyTx,
			b.item.TaskID,
			emptyIfBlank(b.item.Title, "(без названия)"),
			b.item.StatusLabel,
			b.item.RiskScore,
			b.item.HasBlockers,
		))
	}

	lines = append(lines, "")
	lines = append(lines, "Рекомендации SLA:")
	if overdue > 0 {
		lines = append(lines, "- P0: немедленно эскалируйте просроченные задачи и согласуйте новый commit date.")
	}

	if today > 0 {
		lines = append(lines, "- P1: зафиксируйте часовой план на сегодня по задачам с дедлайном <=24ч.")
	}

	if noDeadline > 0 {
		lines = append(lines, "- P2: назначьте дедлайны активным задачам без срока.")
	}

	if overdue == 0 && today == 0 && noDeadline == 0 {
		lines = append(lines, "- SLA в зеленой зоне, продолжайте текущий ритм контроля.")
	}

	lines = append(lines, "")
	lines = append(lines, "=== Вывод ===")
	lines = append(lines, slaConclusion(overdue, today, soon, noDeadline))

	return strings.Join(lines, "\n"), nil
}

func slaConclusion(overdue, today, soon, noDeadline int) string {
	if overdue > 0 {
		return formatConclusion("критично", fmt.Sprintf("SLA нарушен, просрочено %d", overdue), "эскалация и перепланирование просроченных задач", "немедленно")
	}

	if today > 0 {
		return formatConclusion("на грани", fmt.Sprintf("критичных на сегодня: %d", today), "часовой план закрытия задач P1", "сегодня")
	}

	if noDeadline > 0 {
		return formatConclusion("под контролем", fmt.Sprintf("без срока: %d", noDeadline), "назначить дедлайны активным задачам", "в течение 1 рабочего дня")
	}

	if soon > 0 {
		return formatConclusion("желтая зона", fmt.Sprintf("скоро дедлайн: %d", soon), "превентивный контроль и выравнивание нагрузки", "в ближайшие 24 часа")
	}

	return formatConclusion("стабильно", "SLA в зеленой зоне", "сохранить текущий контроль", "по текущему регламенту")
}

func formatConclusion(status, risk, action, reactionTime string) string {
	return strings.Join([]string{
		fmt.Sprintf("Статус: %s", status),
		fmt.Sprintf("Главный риск: %s", risk),
		fmt.Sprintf("Приоритетное действие: %s", action),
		fmt.Sprintf("Срок реакции: %s", reactionTime),
	}, "\n")
}

func runWorkloadSummary(ctx context.Context, client *bitrixClient, filter map[string]any, order map[string]any, start *int, limit *int, includeComments *bool, overloadTasks *int) (string, error) {
	now := time.Now()
	maxTasks := 40
	if limit != nil {
		maxTasks = *limit
	}
	overloadThreshold := 12
	if overloadTasks != nil {
		overloadThreshold = *overloadTasks
	}

	withComments := false
	if includeComments != nil {
		withComments = *includeComments
	}

	tasks, err := loadTaskList(ctx, client, filter, order, start, maxTasks)
	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	items := make([]taskAnalyticsItem, 0, len(tasks))
	for _, task := range tasks {
		id := numberLike(field(task, "id", "ID"))
		if id <= 0 {
			continue
		}

		comments := []map[string]any(nil)
		if withComments {
			comments = loadTaskCommentsSoft(ctx, client, id)
		}
		items = append(items, buildAnalyticsItem(task, comments, now))
	}

	if len(items) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	type workload struct {
		Responsible string
		Active      int
		Overdue     int
		HighRisk    int
		Blockers    int
	}

	m := map[string]*workload{}
	totalActive := 0
	totalOverdue := 0
	totalHighRisk := 0
	for _, it := range items {
		if it.StatusCode == 5 || it.StatusCode == 7 {
			continue
		}

		totalActive++
		if it.IsOverdue {
			totalOverdue++
		}

		if it.RiskScore >= 5 {
			totalHighRisk++
		}

		responsible := strings.TrimSpace(stringField(it.Task, "responsibleId", "RESPONSIBLE_ID"))
		if responsible == "" {
			responsible = "(не указан)"
		}

		w, ok := m[responsible]
		if !ok {
			w = &workload{Responsible: responsible}
			m[responsible] = w
		}

		w.Active++
		if it.IsOverdue {
			w.Overdue++
		}

		if it.RiskScore >= 5 {
			w.HighRisk++
		}

		if it.HasBlockers {
			w.Blockers++
		}
	}

	rows := make([]workload, 0, len(m))
	overloaded := 0
	for _, v := range m {
		rows = append(rows, *v)
		if v.Active >= overloadThreshold {
			overloaded++
		}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Active != rows[j].Active {
			return rows[i].Active > rows[j].Active
		}
		if rows[i].Overdue != rows[j].Overdue {
			return rows[i].Overdue > rows[j].Overdue
		}
		return rows[i].HighRisk > rows[j].HighRisk
	})

	lines := []string{
		fmt.Sprintf("Workload summary (порог перегруза: %d активных задач)", overloadThreshold),
		fmt.Sprintf("Охват: %d задач | Активных: %d | Просрочено: %d | Высокий риск: %d", len(items), totalActive, totalOverdue, totalHighRisk),
		fmt.Sprintf("Ответственных в перегрузе: %d", overloaded),
		"",
		"Нагрузка по ответственным:",
	}
	for _, row := range rows {
		flag := ""
		if row.Active >= overloadThreshold {
			flag = " [ПЕРЕГРУЗ]"
		}

		lines = append(lines, fmt.Sprintf("- %s: активных=%d, просрочено=%d, высокий риск=%d, блокеры=%d%s",
			row.Responsible, row.Active, row.Overdue, row.HighRisk, row.Blockers, flag))
	}

	lines = append(lines, "")
	lines = append(lines, "Рекомендации по перераспределению:")
	lines = append(lines, "- Разгрузить сотрудников с флагом [ПЕРЕГРУЗ] за счет переноса низкоприоритетных задач.")
	lines = append(lines, "- Сначала перераспределять просроченные и high-risk задачи, затем задачи без блокеров.")
	lines = append(lines, "- Зафиксировать owner и дедлайн для каждой задачи, переданной между ответственными.")
	lines = append(lines, "")
	lines = append(lines, "=== Вывод ===")
	if overloaded > 0 || totalOverdue > 0 || totalHighRisk > 0 {
		lines = append(lines, formatConclusion(
			"под нагрузкой",
			fmt.Sprintf("перегруженных: %d, просроченных: %d, high-risk: %d", overloaded, totalOverdue, totalHighRisk),
			"выравнивание нагрузки и приоритизация критичных задач",
			"сегодня",
		))
	} else {
		lines = append(lines, formatConclusion(
			"стабильно",
			"перегруз и критичные риски не выявлены",
			"поддерживать текущий баланс работ",
			"по текущему регламенту",
		))
	}

	return strings.Join(lines, "\n"), nil
}

func runStatusTrends(ctx context.Context, client *bitrixClient, filter map[string]any, order map[string]any, start *int, limit *int, periodDays *int) (string, error) {
	now := time.Now()
	maxTasks := 50
	if limit != nil {
		maxTasks = *limit
	}

	days := 7
	if periodDays != nil {
		days = *periodDays
	}

	currentStart := now.Add(-time.Duration(days) * 24 * time.Hour)
	previousStart := currentStart.Add(-time.Duration(days) * 24 * time.Hour)

	tasks, err := loadTaskList(ctx, client, filter, order, start, maxTasks)
	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "По заданным условиям задачи не найдены.", nil
	}

	current := statusCounters{}
	previous := statusCounters{}
	totalCurrent := 0
	totalPrevious := 0

	for _, t := range tasks {
		status := numberLike(field(t, "status", "STATUS"))
		changedAt := parseBitrixTime(stringField(t, "changedDate", "CHANGED_DATE"))
		if changedAt.IsZero() {
			continue
		}

		if changedAt.After(currentStart) {
			totalCurrent++
			incrementStatusCounter(&current, status)
			continue
		}

		if changedAt.After(previousStart) && changedAt.Before(currentStart) {
			totalPrevious++
			incrementStatusCounter(&previous, status)
		}
	}

	lines := []string{
		fmt.Sprintf("Status trends за %d дн. (по CHANGED_DATE)", days),
		fmt.Sprintf("Текущий период: %d | Предыдущий период: %d", totalCurrent, totalPrevious),
		"",
		"Тренд по корзинам статусов:",
		fmt.Sprintf("- open: %d (%s)", current.Open, trendDelta(current.Open, previous.Open)),
		fmt.Sprintf("- in-progress: %d (%s)", current.InProgress, trendDelta(current.InProgress, previous.InProgress)),
		fmt.Sprintf("- done: %d (%s)", current.Done, trendDelta(current.Done, previous.Done)),
		fmt.Sprintf("- deferred: %d (%s)", current.Deferred, trendDelta(current.Deferred, previous.Deferred)),
		"",
		"=== Вывод ===",
		statusTrendConclusion(current, previous),
	}

	return strings.Join(lines, "\n"), nil
}

func incrementStatusCounter(c *statusCounters, status int) {
	switch status {
	case 3, 4:
		c.InProgress++
	case 5, 7:
		c.Done++
	case 6:
		c.Deferred++
	default:
		c.Open++
	}
}

func statusTrendConclusion(current, previous statusCounters) string {
	criticalRisk := ""
	action := "поддерживать текущий ритм мониторинга статусов"
	reaction := "по текущему регламенту"
	status := "стабильно"

	if current.Deferred > previous.Deferred {
		criticalRisk = fmt.Sprintf("рост deferred (%d -> %d)", previous.Deferred, current.Deferred)
		action = "разобрать причины отложенных задач и вернуть часть в активную работу"
		reaction = "в течение 1 рабочего дня"
		status = "под контролем"
	}

	if current.Open > previous.Open && current.Done <= previous.Done {
		criticalRisk = fmt.Sprintf("рост open без роста done (open %d -> %d, done %d -> %d)", previous.Open, current.Open, previous.Done, current.Done)
		action = "перераспределить фокус команды на закрытие, ограничить входящий поток"
		reaction = "сегодня"
		status = "на грани"
	}

	if current.Done < previous.Done && current.InProgress > previous.InProgress {
		criticalRisk = fmt.Sprintf("снижение закрытия при росте in-progress (done %d -> %d, in-progress %d -> %d)", previous.Done, current.Done, previous.InProgress, current.InProgress)
		action = "сократить WIP и довести задачи до завершения"
		reaction = "сегодня"
		status = "критично"
	}

	if criticalRisk == "" {
		criticalRisk = "негативный тренд по статусам не выявлен"
	}

	return formatConclusion(status, criticalRisk, action, reaction)
}
