package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type assistantMessageRegenerationRepository struct {
	db *pgxpool.Pool
}

func NewAssistantMessageRegenerationRepository(db *pgxpool.Pool) domain.AssistantMessageRegenerationRepository {
	return &assistantMessageRegenerationRepository{db: db}
}

func (r *assistantMessageRegenerationRepository) Create(ctx context.Context, regen *domain.AssistantMessageRegeneration) error {
	if regen == nil {
		return nil
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO message_edits
			(session_id, message_id, editor_user_id, kind, old_content, new_content, soft_deleted_from_id, soft_deleted_to_id, created_at)
		VALUES ($1, $2, $3, 'assistant_regen', $4, $5, NULL, NULL, $6)
		RETURNING id
	`,
		regen.SessionId,
		regen.MessageId,
		regen.RegenUserId,
		regen.OldContent,
		regen.NewContent,
		regen.CreatedAt,
	).Scan(&regen.Id)
}

func (r *assistantMessageRegenerationRepository) ListByMessageID(ctx context.Context, messageID int64, limit int32) ([]*domain.AssistantMessageRegeneration, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, session_id, message_id, editor_user_id, old_content, new_content, created_at
		FROM message_edits
		WHERE message_id = $1 AND kind = 'assistant_regen'
		ORDER BY created_at DESC, id DESC
		LIMIT $2
	`, messageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*domain.AssistantMessageRegeneration, 0, limit)
	for rows.Next() {
		var row domain.AssistantMessageRegeneration
		if err := rows.Scan(
			&row.Id,
			&row.SessionId,
			&row.MessageId,
			&row.RegenUserId,
			&row.OldContent,
			&row.NewContent,
			&row.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, &row)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return out, nil
}
