package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/log"
)

// Seed is the interface to implement to provide seeds.
type Seed interface {
	Name() string
	Run(ctx context.Context, db *sql.DB) error
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
		seedErr := seedToExecute.Run(ctx, m.db)
		if seedErr != nil {
			m.logger.Error().Err(seedErr).Str("seed", seedToExecute.Name()).Msg("seed error")

			return seedErr
		}

		m.logger.Info().Str("seed", seedToExecute.Name()).Msg("seed success")
	}

	return nil
}
