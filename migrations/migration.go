package migrations

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	sqlDB := stdlib.OpenDBFromPool(db)
	defer sqlDB.Close()
	return goose.UpContext(ctx, sqlDB, ".")
}
