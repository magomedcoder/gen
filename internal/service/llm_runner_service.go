package service

import (
	"context"
	"fmt"
	"github.com/magomedcoder/gen/api/pb/llmrunner"
	"strings"

	"github.com/magomedcoder/gen/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LLMRunnerService struct {
	client llmrunnerpb.LLMRunnerServiceClient
	conn   *grpc.ClientConn
	model  string
}

func NewLLMRunnerService(address, model string) (*LLMRunnerService, error) {
	if address == "" {
		address = "localhost:50052"
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("подключение к llm-runner: %w", err)
	}
	return &LLMRunnerService{
		client: llmrunnerpb.NewLLMRunnerServiceClient(conn),
		conn:   conn,
		model:  model,
	}, nil
}

func (s *LLMRunnerService) Close() error {
	return s.conn.Close()
}

func (s *LLMRunnerService) CheckConnection(ctx context.Context) (bool, error) {
	resp, err := s.client.CheckConnection(ctx, &llmrunnerpb.Empty{})
	if err != nil {
		return false, fmt.Errorf("llm-runner CheckConnection: %w", err)
	}
	return resp.IsConnected, nil
}

func (s *LLMRunnerService) GetModels(ctx context.Context) ([]string, error) {
	resp, err := s.client.GetModels(ctx, &llmrunnerpb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("llm-runner GetModels: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	return resp.Models, nil
}

func (s *LLMRunnerService) GetGpuInfo(ctx context.Context) (*llmrunnerpb.GetGpuInfoResponse, error) {
	resp, err := s.client.GetGpuInfo(ctx, &llmrunnerpb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("llm-runner GetGpuInfo: %w", err)
	}

	return resp, nil
}

func (s *LLMRunnerService) GetServerInfo(ctx context.Context) (*llmrunnerpb.ServerInfo, error) {
	resp, err := s.client.GetServerInfo(ctx, &llmrunnerpb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("llm-runner GetServerInfo: %w", err)
	}

	return resp, nil
}

func (s *LLMRunnerService) SendMessage(ctx context.Context, sessionID int64, model string, messages []*domain.Message) (chan string, error) {
	modelName := strings.TrimSpace(model)
	if modelName == "" || modelName == "default" {
		modelName = strings.TrimSpace(s.model)
	}
	if modelName == "default" {
		modelName = ""
	}
	req := &llmrunnerpb.SendMessageRequest{
		SessionId: sessionID,
		Messages:  domainMessagesToProto(messages),
		Model:     modelName,
	}

	stream, err := s.client.SendMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("llm-runner SendMessage: %w", err)
	}

	output := make(chan string, 100)

	go func() {
		defer close(output)
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			if msg.Content != "" {
				select {
				case <-ctx.Done():
					return
				case output <- msg.Content:
				}
			}
			if msg.Done {
				return
			}
		}
	}()

	return output, nil
}

func domainMessagesToProto(messages []*domain.Message) []*llmrunnerpb.ChatMessage {
	out := make([]*llmrunnerpb.ChatMessage, len(messages))
	for i, m := range messages {
		out[i] = &llmrunnerpb.ChatMessage{
			Id:        int64(i + 1),
			Content:   m.Content,
			Role:      string(m.Role),
			CreatedAt: m.CreatedAt.Unix(),
		}
	}
	return out
}
