package handler

import (
	"context"

	"github.com/magomedcoder/gen/api/pb/runnerpb"
	"github.com/magomedcoder/gen/internal/runner"
	"github.com/magomedcoder/gen/internal/usecase"
)

type RunnerHandler struct {
	runnerpb.UnimplementedRunnerServiceServer
	runnerpb.UnimplementedRunnerAdminServiceServer
	registry    *runner.Registry
	authUseCase *usecase.AuthUseCase
}

func NewRunnerHandler(registry *runner.Registry, authUseCase *usecase.AuthUseCase) *RunnerHandler {
	return &RunnerHandler{
		registry:    registry,
		authUseCase: authUseCase,
	}
}

func (h *RunnerHandler) GetRunners(ctx context.Context, _ *runnerpb.Empty) (*runnerpb.GetRunnersResponse, error) {
	if err := RequireAdmin(ctx, h.authUseCase); err != nil {
		return nil, err
	}

	return &runnerpb.GetRunnersResponse{
		Runners: h.registry.GetRunners(),
	}, nil
}

func (h *RunnerHandler) SetRunnerEnabled(ctx context.Context, req *runnerpb.SetRunnerEnabledRequest) (*runnerpb.Empty, error) {
	if err := RequireAdmin(ctx, h.authUseCase); err != nil {
		return nil, err
	}
	if req != nil {
		h.registry.SetEnabled(req.Address, req.Enabled)
	}

	return &runnerpb.Empty{}, nil
}

func (h *RunnerHandler) GetRunnersStatus(ctx context.Context, _ *runnerpb.Empty) (*runnerpb.GetRunnersStatusResponse, error) {
	return &runnerpb.GetRunnersStatusResponse{
		HasActiveRunners: h.registry.HasActiveRunners(),
	}, nil
}
