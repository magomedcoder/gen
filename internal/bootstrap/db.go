package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/magomedcoder/gen/config"
)

func CheckDatabase(ctx context.Context, dbCfg config.DatabaseConfig) error {
	targetDB, err := dbCfg.TargetDBName()
	if err != nil {
		return fmt.Errorf("конфигурация базы данных: %w", err)
	}
	adminDSN, err := dbCfg.AdminPostgresDSN()
	if err != nil {
		return fmt.Errorf("конфигурация базы данных (admin): %w", err)
	}

	pool, err := pgxpool.New(ctx, adminDSN)
	if err != nil {
		return fmt.Errorf("ошибка подключения к postgres: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ошибка проверки соединения с postgres: %w", err)
	}

	var exists int
	if err = pool.QueryRow(ctx, "SELECT 1 FROM pg_database WHERE datname = $1", targetDB).Scan(&exists); err == nil {
		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("ошибка проверки существования БД: %w", err)
	}

	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", quoteIdentifier(targetDB)))
	if err != nil {
		return fmt.Errorf("ошибка создания базы данных %s: %w", targetDB, err)
	}

	return nil
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
