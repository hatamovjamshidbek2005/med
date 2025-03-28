package services

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"med/internal/db/psql/sqlc"
	"med/pkg/email"
	"med/pkg/logger"
)

type Service struct {
	storage sqlc.Querier
	log     logger.ILogger
	pool    *pgxpool.Pool
	email   email.Email
}

func NewService(storage sqlc.Querier, pool *pgxpool.Pool, log logger.ILogger, email email.Email) *Service {
	service := &Service{
		storage: storage,
		log:     log,
		pool:    pool,
		email:   email,
	}

	return service
}
