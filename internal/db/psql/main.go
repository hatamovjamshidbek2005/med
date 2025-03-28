package psql

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"med/internal/configs"
	"med/internal/db/psql/sqlc"
	"med/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	DB  *pgxpool.Pool
	log logger.ILogger
	cfg configs.Config
	*sqlc.Queries
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, log logger.ILogger, cnf configs.Config) (*Store, error) {
	pool, err := pgxpool.Connect(ctx, cnf.DBSource)
	if err != nil {
		log.Error("Error creating connection pool", logger.Error(err))
		return nil, err
	}

	//migrationPath := "file://internal/db/psql/migration"
	//m, err := migrate.New(migrationPath, cnf.DBSource)
	//if err != nil {
	//	log.Error("Failed to initialize migration", logger.Error(err))
	//	pool.Close()
	//	return nil, err
	//}

	//if err = m.Up(); err != nil {
	//	if !strings.Contains(err.Error(), "no change") {
	//		log.Warn("Migration failed", logger.Error(err))
	//		version, dirty, verr := m.Version()
	//		if verr != nil {
	//			log.Error("Failed to get migration version", logger.Error(verr))
	//			pool.Close()
	//			return nil, verr
	//		}
	//		if dirty {
	//			log.Warn("Migration is dirty, forcing previous version", logger.Any("version", version))
	//			version--
	//			if err = m.Force(int(version)); err != nil {
	//				log.Error("Failed to force migration version", logger.Error(err))
	//				pool.Close()
	//				return nil, err
	//			}
	//		}
	//		pool.Close()
	//		return nil, err
	//	}
	//	log.Info("No new migrations to apply")
	//} else {
	//	log.Info("Migrations applied successfully")
	//}

	queries := sqlc.New(pool)
	if queries == nil {
		log.Error("Failed to initialize sqlc Queries")
		pool.Close()
		return nil, fmt.Errorf("sqlc Queries is nil")
	}

	store := &Store{
		DB:      pool,
		log:     log,
		cfg:     cnf,
		pool:    pool,
		Queries: queries,
	}

	log.Info("Database store initialized successfully")
	return store, nil
}

func (s *Store) Pool() *pgxpool.Pool {
	return s.pool
}
