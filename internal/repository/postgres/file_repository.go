package postgres

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type fileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) domain.FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *domain.File) error {
	kind := strings.TrimSpace(file.Kind)
	return r.db.QueryRow(ctx, `
		INSERT INTO files (filename, mime_type, size, storage_path, created_at, chat_session_id, user_id, expires_at, kind)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		file.Filename,
		nullStringPtr(file.MimeType),
		file.Size,
		file.StoragePath,
		file.CreatedAt,
		file.ChatSessionID,
		file.UserID,
		file.ExpiresAt,
		kind,
	).Scan(&file.Id)
}

func (r *fileRepository) UpdateStoragePath(ctx context.Context, id int64, storagePath string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE files SET storage_path = $2 WHERE id = $1
	`, id, storagePath)
	return err
}

func (r *fileRepository) GetById(ctx context.Context, id int64) (*domain.File, error) {
	var f domain.File
	var mimeType *string
	var chatSessionID *int64
	var userID *int
	var expiresAt *time.Time
	err := r.db.QueryRow(ctx, `
		SELECT id, filename, mime_type, size, storage_path, created_at, chat_session_id, user_id, expires_at, kind
		FROM files
		WHERE id = $1
	`, id).Scan(
		&f.Id,
		&f.Filename,
		&mimeType,
		&f.Size,
		&f.StoragePath,
		&f.CreatedAt,
		&chatSessionID,
		&userID,
		&expiresAt,
		&f.Kind,
	)
	if err != nil {
		return nil, err
	}

	if mimeType != nil {
		f.MimeType = *mimeType
	}
	f.ChatSessionID = chatSessionID
	f.UserID = userID
	f.ExpiresAt = expiresAt

	return &f, nil
}

func (r *fileRepository) DeleteExpired(ctx context.Context) (deleted int64, err error) {
	rows, err := r.db.Query(ctx, `
		DELETE FROM files
		WHERE expires_at IS NOT NULL AND expires_at < NOW()
		RETURNING storage_path
	`)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	var n int64
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return n, err
		}

		if p != "" && p != "." {
			_ = os.Remove(p)
		}

		n++
	}

	return n, rows.Err()
}

func (r *fileRepository) CountSessionTTLArtifacts(ctx context.Context, sessionID int64, userID int) (count int32, totalSize int64, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)::int, COALESCE(SUM(size), 0)::bigint
		FROM files
		WHERE chat_session_id = $1 AND user_id = $2 AND kind = 'artifact'
		  AND expires_at IS NOT NULL AND expires_at > NOW()
	`, sessionID, userID).Scan(&count, &totalSize)
	return count, totalSize, err
}

func nullStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
