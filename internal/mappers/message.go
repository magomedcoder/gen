package mappers

import (
	"time"

	"github.com/magomedcoder/gen/api/pb"
	"github.com/magomedcoder/gen/internal/domain"
)

func MessageToProto(msg *domain.Message) *pb.ChatMessage {
	if msg == nil {
		return nil
	}

	return &pb.ChatMessage{
		Id:        msg.Id,
		Content:   msg.Content,
		Role:      domain.ToProtoRole(msg.Role),
		CreatedAt: msg.CreatedAt.Unix(),
	}
}

func MessagesFromProto(pbMsgs []*pb.ChatMessage, sessionID string) []*domain.Message {
	if len(pbMsgs) == 0 {
		return nil
	}
	out := make([]*domain.Message, 0, len(pbMsgs))
	for _, m := range pbMsgs {
		if m == nil {
			continue
		}
		var createdAt time.Time
		if m.CreatedAt != 0 {
			createdAt = time.Unix(m.CreatedAt, 0)
		}
		out = append(out, &domain.Message{
			Id:        m.Id,
			SessionId: sessionID,
			Content:   m.Content,
			Role:      domain.FromProtoRole(m.Role),
			CreatedAt: createdAt,
		})
	}
	return out
}
