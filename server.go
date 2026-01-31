package runner

import (
	"context"

	"github.com/magomedcoder/llm-runner/pb"
)

type Server struct {
	pb.UnimplementedLLMRunnerServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Ok: false,
	}, nil
}

func (s *Server) GetModels(ctx context.Context, _ *pb.Empty) (*pb.GetModelsResponse, error) {
	return &pb.GetModelsResponse{
		Models: nil,
	}, nil
}
