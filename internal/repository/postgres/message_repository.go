package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type messageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &messageRepository{db: db}
}

func nullInt64Ptr(v *int64) interface{} {
	if v == nil {
		return nil
	}

	return *v
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	err := r.db.QueryRow(ctx, `
		INSERT INTO messages (session_id, content, role, attachment_file_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		message.SessionId,
		message.Content,
		message.Role,
		nullInt64Ptr(message.AttachmentFileID),
		message.CreatedAt,
		message.UpdatedAt,
	).Scan(&message.Id)

	return err
}

func (r *messageRepository) UpdateContent(ctx context.Context, id int64, content string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE messages
		SET content = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id, content)
	return err
}

func (r *messageRepository) GetBySessionId(ctx context.Context, sessionID int64, page, pageSize int32) ([]*domain.Message, int32, error) {
	_, pageSize, offset := normalizePagination(page, pageSize)

	var total int32
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM messages
		WHERE session_id = $1 AND deleted_at IS NULL
	`, sessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT m.id, m.session_id, m.content, m.role, m.attachment_file_id, f.filename, m.created_at, m.updated_at, m.deleted_at
		FROM messages m
		LEFT JOIN files f ON m.attachment_file_id = f.id
		WHERE m.session_id = $1 AND m.deleted_at IS NULL
		ORDER BY m.created_at ASC
		LIMIT $2 OFFSET $3
	`, sessionID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		var attachmentFileID *int64
		var fname *string
		if err := rows.Scan(
			&message.Id,
			&message.SessionId,
			&message.Content,
			&message.Role,
			&attachmentFileID,
			&fname,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		message.AttachmentFileID = attachmentFileID
		if fname != nil {
			message.AttachmentName = *fname
		}
		messages = append(messages, &message)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return messages, total, nil
}
