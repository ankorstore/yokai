package seed

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
)

type ValidSeed struct {
	config *config.Config
}

func NewValidSeed(config *config.Config) *ValidSeed {
	return &ValidSeed{
		config: config,
	}
}

func (s *ValidSeed) Name() string {
	return "valid"
}

func (s *ValidSeed) Run(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO foo (bar) VALUES (?)", s.config.GetString("config.seed_value"))

	return err
}
