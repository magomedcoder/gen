package bitrix24server

import "testing"

func pint(v int) *int { return &v }

func TestValidateAnalyzeToolsArgs(t *testing.T) {
	tests := []struct {
		name    string
		run     func() error
		wantErr bool
	}{
		{"list-valid", func() error {
			return validateListTasksArgs(pint(0))
		}, false},

		{"list-invalid", func() error {
			return validateListTasksArgs(pint(-1))
		}, true},

		{"query-valid", func() error {
			return validateAnalyzeTasksByQueryArgs(pint(0), pint(20))
		}, false},

		{"query-invalid-limit", func() error {
			return validateAnalyzeTasksByQueryArgs(pint(0), pint(51))
		}, true},

		{"portfolio-valid", func() error {
			return validatePortfolioArgs(pint(0), pint(30), "responsible")
		}, false},

		{"portfolio-invalid-group", func() error {
			return validatePortfolioArgs(pint(0), pint(30), "team")
		}, true},

		{"exec-valid", func() error {
			return validateExecutiveSummaryArgs(pint(0), pint(40), pint(7))
		}, false},

		{"exec-invalid-period", func() error {
			return validateExecutiveSummaryArgs(pint(0), pint(40), pint(31))
		}, true},

		{"sla-valid", func() error {
			return validateSLAArgs(pint(0), pint(40), pint(24))
		}, false},

		{"sla-invalid-threshold", func() error {
			return validateSLAArgs(pint(0), pint(40), pint(169))
		}, true},

		{"workload-valid", func() error {
			return validateWorkloadArgs(pint(0), pint(40), pint(12))
		}, false},

		{"workload-invalid-overload", func() error {
			return validateWorkloadArgs(pint(0), pint(40), pint(101))
		}, true},

		{"status-valid", func() error {
			return validateStatusTrendsArgs(pint(0), pint(50), pint(7))
		}, false},

		{"status-invalid-limit", func() error {
			return validateStatusTrendsArgs(pint(0), pint(0), pint(7))
		}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
