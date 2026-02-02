package runner

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/magomedcoder/llm-runner/domain"
	"github.com/magomedcoder/llm-runner/pb"
	"github.com/magomedcoder/llm-runner/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedLLMRunnerServiceServer
	textProvider     provider.TextProvider
	inferenceMetrics *InferenceMetrics
	sem              chan struct{}
	addresses        []string
	addressesMu      sync.Mutex
}

func NewServer(textProvider provider.TextProvider, maxConcurrentGenerations int) *Server {
	var sem chan struct{}
	if maxConcurrentGenerations > 0 {
		sem = make(chan struct{}, maxConcurrentGenerations)
	}

	return &Server{
		textProvider:     textProvider,
		inferenceMetrics: NewInferenceMetrics(),
		sem:              sem,
	}
}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	if s.textProvider == nil {
		return &pb.PingResponse{Ok: false}, nil
	}

	ok, _ := s.textProvider.CheckConnection(ctx)

	return &pb.PingResponse{Ok: ok}, nil
}

func (s *Server) GetModels(ctx context.Context, _ *pb.Empty) (*pb.GetModelsResponse, error) {
	if s.textProvider == nil {
		return &pb.GetModelsResponse{}, nil
	}

	models, err := s.textProvider.GetModels(ctx)
	if err != nil {
		return &pb.GetModelsResponse{}, nil
	}

	return &pb.GetModelsResponse{
		Models: models,
	}, nil
}

func (s *Server) Generate(req *pb.GenerateRequest, stream pb.LLMRunnerService_GenerateServer) error {
	if s.textProvider == nil {
		return status.Error(codes.Unavailable, "поставщик текста не задан")
	}

	if req == nil || len(req.Messages) == 0 {
		return stream.Send(&pb.GenerateResponse{Done: true})
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

	sessionId := req.SessionId
	model := req.Model
	messages := domain.AIMessagesFromProto(req.Messages, sessionId)
	stopSequences := req.GetStopSequences()
	genParams := buildGenParamsFromRequest(req)
	if s := req.GetTimeoutSeconds(); s > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(s)*time.Second)
		defer cancel()
	}

	start := time.Now()
	var tokens int64
	var fullContent strings.Builder
	ch, err := s.textProvider.SendMessage(ctx, sessionId, model, messages, stopSequences, genParams)
	if err != nil {
		_ = stream.Send(&pb.GenerateResponse{Done: true})
		return err
	}

	for chunk := range ch {
		if chunk != "" {
			tokens++
			fullContent.WriteString(chunk)
			if err := stream.Send(&pb.GenerateResponse{Content: chunk, Done: false}); err != nil {
				return err
			}
		}
	}

	if s.inferenceMetrics != nil {
		s.inferenceMetrics.Record(tokens, time.Since(start))
	}

	resp := &pb.GenerateResponse{Done: true}
	if len(req.Tools) > 0 {
		if toolCalls := ParseToolCalls(fullContent.String()); len(toolCalls) > 0 {
			resp.ToolCalls = make([]*pb.ToolCall, len(toolCalls))
			for i, tc := range toolCalls {
				resp.ToolCalls[i] = &pb.ToolCall{Id: tc.Id, Name: tc.Name, Arguments: tc.Arguments}
			}
		}
	}

	return stream.Send(resp)
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRunnerRequest) (*pb.Empty, error) {
	if req != nil && req.Address != "" {
		s.addressesMu.Lock()
		s.addresses = append(s.addresses, req.Address)
		s.addressesMu.Unlock()
	}

	return &pb.Empty{}, nil
}

func (s *Server) Unregister(ctx context.Context, req *pb.UnregisterRunnerRequest) (*pb.Empty, error) {
	if req != nil && req.Address != "" {
		s.addressesMu.Lock()
		for i, a := range s.addresses {
			if a == req.Address {
				s.addresses = append(s.addresses[:i], s.addresses[i+1:]...)
				break
			}
		}

		s.addressesMu.Unlock()
	}

	return &pb.Empty{}, nil
}

func buildGenParamsFromRequest(req *pb.GenerateRequest) *domain.GenerationParams {
	if req == nil {
		return nil
	}

	hasSampling := req.Temperature != nil || req.MaxTokens != nil || req.TopK != nil || req.TopP != nil
	hasFormat := req.ResponseFormat != nil
	hasTools := len(req.Tools) > 0
	if !hasSampling && !hasFormat && !hasTools {
		return nil
	}

	p := &domain.GenerationParams{
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopK:        req.TopK,
		TopP:        req.TopP,
	}

	if hasFormat {
		p.ResponseFormat = &domain.ResponseFormat{
			Type:   req.ResponseFormat.Type,
			Schema: req.ResponseFormat.Schema,
		}
	}

	if hasTools {
		p.Tools = make([]domain.Tool, len(req.Tools))
		for i, t := range req.Tools {
			p.Tools[i] = domain.Tool{
				Name:           t.Name,
				Description:    t.Description,
				ParametersJSON: t.ParametersJson,
			}
		}
	}

	return p
}
