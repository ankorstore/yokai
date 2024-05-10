package fxsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ankorstore/yokai/log"
)

// Seed is the interface to implement to provide seeds.
type Seed interface {
	Name() string
	Run(ctx context.Context, tx *sql.Tx) error
}

// Seeder is a database seeder.
type Seeder struct {
	db     *sql.DB
	logger *log.Logger
	seeds  []Seed
}

// NewSeeder returns a new Seeder instance.
func NewSeeder(db *sql.DB, logger *log.Logger, seeds ...Seed) *Seeder {
	return &Seeder{
		db:     db,
		logger: logger,
		seeds:  seeds,
	}
}

// Run executes the Seed list matching to the provided list of names, or all of them if the list of names is empty.
func (m *Seeder) Run(ctx context.Context, names ...string) error {
	var seedsToExecute []Seed

	if len(names) == 0 {
		seedsToExecute = m.seeds
	} else {
		for _, name := range names {
			for _, seed := range m.seeds {
				if name == seed.Name() {
					seedsToExecute = append(seedsToExecute, seed)
				}
			}
		}
	}

	for _, seedToExecute := range seedsToExecute {
		tx, err := m.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("cannot begin transaction for seed %s: %w", seedToExecute.Name(), err)
		}

		seedErr := seedToExecute.Run(ctx, tx)
		if seedErr != nil {
			err = tx.Rollback()
			if err != nil {
				return fmt.Errorf("cannot rollback transaction for seed %s: %w", seedToExecute.Name(), err)
			}

			m.logger.Error().Err(seedErr).Str("seed", seedToExecute.Name()).Msg("rollback")
		} else {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("cannot commit transaction for seed %s: %w", seedToExecute.Name(), err)
			}

			m.logger.Info().Str("seed", seedToExecute.Name()).Msg("commit")
		}
	}

	return nil
}
