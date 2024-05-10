package seed

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
)

type InvalidSeed struct {
	config *config.Config
}

func NewInvalidSeed(config *config.Config) *InvalidSeed {
	return &InvalidSeed{
		config: config,
	}
}

func (s *InvalidSeed) Name() string {
	return "invalid"
}

func (s *InvalidSeed) Run(ctx context.Context, tx *sql.Tx) error {
	// should succeed
	_, err := tx.ExecContext(ctx, "INSERT INTO foo (bar) VALUES (?)", s.config.GetString("config.seed_value"))
	if err != nil {
		return err
	}

	// should fail
	_, err = tx.ExecContext(ctx, "INSERT INTO invalid (bar) VALUES (?)", s.config.GetString("config.seed_value"))
	if err != nil {
		return err
	}

	return nil
}
