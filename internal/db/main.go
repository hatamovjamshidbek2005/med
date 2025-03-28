package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"med/internal/configs"
	"med/internal/db/psql"
	"med/internal/db/psql/sqlc"
	"med/pkg/logger"
)

type StorageI interface {
	sqlc.Querier
	Pool() *pgxpool.Pool
}

func New(ctx context.Context, cfg configs.Config, log logger.ILogger) (StorageI, error) {
	return psql.NewStore(ctx, log, cfg)
}
