package seed

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
)

type TestSeed struct {
	config *config.Config
}

func NewTestSeed(config *config.Config) *TestSeed {
	return &TestSeed{
		config: config,
	}
}

func (s *TestSeed) Name() string {
	return "test"
}

func (s *TestSeed) Run(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "INSERT INTO foo (bar) VALUES (?)", s.config.GetString("config.seed_value"))

	return err
}
