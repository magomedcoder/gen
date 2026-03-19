package runner

import (
	"context"
	"github.com/magomedcoder/llm-runner/domain"
	"github.com/magomedcoder/llm-runner/gpu"
	"github.com/magomedcoder/llm-runner/pb/llmrunnerpb"
	"github.com/magomedcoder/llm-runner/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type Server struct {
	llmrunnerpb.UnimplementedLLMRunnerServiceServer
	textProvider     provider.TextProvider
	gpuCollector     gpu.Collector
	inferenceMetrics *InferenceMetrics
	sem              chan struct{}
	defaultModel     string
}

func NewServer(textProvider provider.TextProvider, gpuCollector gpu.Collector, maxConcurrentGenerations int, defaultModel string) *Server {
	if gpuCollector == nil {
		gpuCollector = gpu.NewCollector()
	}
	var sem chan struct{}
	if maxConcurrentGenerations > 0 {
		sem = make(chan struct{}, maxConcurrentGenerations)
	}
	return &Server{
		textProvider:     textProvider,
		gpuCollector:     gpuCollector,
		inferenceMetrics: NewInferenceMetrics(),
		sem:              sem,
		defaultModel:     strings.TrimSpace(defaultModel),
	}
}

func (s *Server) CheckConnection(ctx context.Context, _ *llmrunnerpb.Empty) (*llmrunnerpb.ConnectionResponse, error) {
	if s.textProvider == nil {
		return &llmrunnerpb.ConnectionResponse{IsConnected: false}, nil
	}

	ok, _ := s.textProvider.CheckConnection(ctx)
	return &llmrunnerpb.ConnectionResponse{IsConnected: ok}, nil
}

func (s *Server) GetModels(ctx context.Context, _ *llmrunnerpb.Empty) (*llmrunnerpb.GetModelsResponse, error) {
	if s.textProvider == nil {
		return &llmrunnerpb.GetModelsResponse{}, nil
	}

	models, err := s.textProvider.GetModels(ctx)
	if err != nil {
		return &llmrunnerpb.GetModelsResponse{}, nil
	}

	return &llmrunnerpb.GetModelsResponse{
		Models: models,
	}, nil
}

func (s *Server) SendMessage(req *llmrunnerpb.SendMessageRequest, stream llmrunnerpb.LLMRunnerService_SendMessageServer) error {
	if s.textProvider == nil {
		return status.Error(codes.Unavailable, "текстовый провайдер не подключён")
	}

	if req == nil || len(req.Messages) == 0 {
		return stream.Send(&llmrunnerpb.ChatResponse{Done: true})
	}

	ctx := stream.Context()

	if s.sem != nil {
		select {
		case s.sem <- struct{}{}:
			defer func() { <-s.sem }()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	sessionID := req.SessionId
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = s.defaultModel
	}
	messages := domain.AIMessagesFromProto(req.Messages, sessionID)
	stopSequences := req.GetStopSequences()

	if ts := req.GetTimeoutSeconds(); ts > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(ts)*time.Second)
		defer cancel()
	}

	start := time.Now()
	var tokens int64

	ch, err := s.textProvider.SendMessage(ctx, sessionID, model, messages, stopSequences, nil)
	if err != nil {
		_ = stream.Send(&llmrunnerpb.ChatResponse{Done: true})
		return err
	}

	for chunk := range ch {
		if chunk != "" {
			tokens++
			if err := stream.Send(&llmrunnerpb.ChatResponse{
				Content: chunk,
				Done:    false,
			}); err != nil {
				return err
			}
		}
	}

	if s.inferenceMetrics != nil {
		s.inferenceMetrics.Record(tokens, time.Since(start))
	}

	return stream.Send(&llmrunnerpb.ChatResponse{Done: true})
}

func (s *Server) GetGpuInfo(ctx context.Context, _ *llmrunnerpb.Empty) (*llmrunnerpb.GetGpuInfoResponse, error) {
	list := s.gpuCollector.Collect()
	gpus := make([]*llmrunnerpb.GpuInfo, len(list))
	for i := range list {
		gpus[i] = &llmrunnerpb.GpuInfo{
			Name:               list[i].Name,
			TemperatureC:       list[i].TemperatureC,
			MemoryTotalMb:      list[i].MemoryTotalMB,
			MemoryUsedMb:       list[i].MemoryUsedMB,
			UtilizationPercent: list[i].UtilizationPercent,
		}
	}

	return &llmrunnerpb.GetGpuInfoResponse{Gpus: gpus}, nil
}

func (s *Server) GetServerInfo(ctx context.Context, _ *llmrunnerpb.Empty) (*llmrunnerpb.ServerInfo, error) {
	si := CollectSysInfo()
	out := &llmrunnerpb.ServerInfo{
		Hostname:      si.Hostname,
		Os:            si.OS,
		Arch:          si.Arch,
		CpuCores:      si.CPUCores,
		MemoryTotalMb: si.MemoryTotalMB,
	}
	if s.textProvider != nil {
		if models, err := s.textProvider.GetModels(ctx); err == nil && len(models) > 0 {
			out.Models = models
		}
	}

	return out, nil
}
