package bitrix24server

import "time"

type RiskScoringConfig struct {
	OverdueWeight            int
	NoDeadlineActiveWeight   int
	NoCommentsStaleWeight    int
	NoFreshCommentsWeight    int
	BlockersWeight           int
	EstimateOverrunWeight    int
	NoCommentsStaleAfter     time.Duration
	NoFreshCommentsAfter     time.Duration
	MediumRiskScoreThreshold int
	HighRiskScoreThreshold   int
}

type RiskInput struct {
	Now           time.Time
	StatusCode    int
	CreatedAt     time.Time
	Deadline      time.Time
	LastComment   time.Time
	CommentsCount int
	HasBlockers   bool
	TimeEstimate  int
	TimeSpent     int
}

func defaultRiskScoringConfig() RiskScoringConfig {
	return RiskScoringConfig{
		OverdueWeight:            4,
		NoDeadlineActiveWeight:   2,
		NoCommentsStaleWeight:    2,
		NoFreshCommentsWeight:    2,
		BlockersWeight:           2,
		EstimateOverrunWeight:    2,
		NoCommentsStaleAfter:     48 * time.Hour,
		NoFreshCommentsAfter:     72 * time.Hour,
		MediumRiskScoreThreshold: 2,
		HighRiskScoreThreshold:   5,
	}
}

func evaluateRisk(input RiskInput, cfg RiskScoringConfig) (score int, riskLabel string, reasons []string) {
	if input.Now.IsZero() {
		input.Now = time.Now()
	}

	reasons = make([]string, 0, 6)

	if !input.Deadline.IsZero() && input.Deadline.Before(input.Now) {
		score += cfg.OverdueWeight
		reasons = append(reasons, "задача уже просрочена")
	}

	isActive := input.StatusCode == 2 || input.StatusCode == 3 || input.StatusCode == 4
	if input.Deadline.IsZero() && isActive {
		score += cfg.NoDeadlineActiveWeight
		reasons = append(reasons, "нет дедлайна у активной задачи")
	}

	if input.CommentsCount == 0 && !input.CreatedAt.IsZero() && input.Now.Sub(input.CreatedAt) > cfg.NoCommentsStaleAfter {
		score += cfg.NoCommentsStaleWeight
		reasons = append(reasons, "долгое время нет комментариев")
	}

	if !input.LastComment.IsZero() && input.Now.Sub(input.LastComment) > cfg.NoFreshCommentsAfter && input.StatusCode != 5 && input.StatusCode != 7 {
		score += cfg.NoFreshCommentsWeight
		reasons = append(reasons, "нет свежих коммуникаций по задаче")
	}

	if input.HasBlockers {
		score += cfg.BlockersWeight
		reasons = append(reasons, "в комментариях есть сигналы блокировки/ожидания")
	}

	if input.TimeEstimate > 0 && input.TimeSpent > input.TimeEstimate {
		score += cfg.EstimateOverrunWeight
		reasons = append(reasons, "превышена оценка времени")
	}

	if input.StatusCode == 5 || input.StatusCode == 7 {
		score = 0
		reasons = nil
	}

	switch {
	case score >= cfg.HighRiskScoreThreshold:
		riskLabel = "высокий"
	case score >= cfg.MediumRiskScoreThreshold:
		riskLabel = "средний"
	default:
		riskLabel = "низкий"
	}

	return score, riskLabel, reasons
}
