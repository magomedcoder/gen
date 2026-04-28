package bitrix24server

import (
	"context"
	"time"
)

type AnalyticsContext struct {
	Now   time.Time
	Items []taskAnalyticsItem
}

func buildAnalyticsContextForTaskList(ctx context.Context, client *bitrixClient, filter map[string]any, order map[string]any, start *int, limit int, includeComments bool) (*AnalyticsContext, error) {
	now := time.Now()
	tasks, err := loadTaskList(ctx, client, filter, order, start, limit)
	if err != nil {
		return nil, err
	}

	items := make([]taskAnalyticsItem, 0, len(tasks))
	for _, task := range tasks {
		id := numberLike(field(task, "id", "ID"))
		if id <= 0 {
			continue
		}

		comments := []map[string]any(nil)
		if includeComments {
			comments = loadTaskCommentsSoft(ctx, client, id)
		}

		items = append(items, buildAnalyticsItem(task, comments, now))
	}

	return &AnalyticsContext{
		Now:   now,
		Items: items,
	}, nil
}
