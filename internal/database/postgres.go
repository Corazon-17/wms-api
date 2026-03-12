package database

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"wms-api/internal/config"
)

func NewDB(cfg *config.Config) *bun.DB {

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.DatabaseURL),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	db.PingContext(context.Background())

	return db
}
