package healthcheck

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log"
)

// DefaultProbeName is the name of the SQL probe.
const DefaultProbeName = "sql"

// SQLProbe is a probe compatible with the [healthcheck] module.
//
// [healthcheck]: https://github.com/ankorstore/yokai/tree/main/healthcheck
type SQLProbe struct {
	name string
	db   *sql.DB
}

// NewSQLProbe returns a new [SQLProbe].
func NewSQLProbe(db *sql.DB) *SQLProbe {
	return &SQLProbe{
		name: DefaultProbeName,
		db:   db,
	}
}

// Name returns the name of the [SQLProbe].
func (p *SQLProbe) Name() string {
	return p.name
}

// SetName sets the name of the [SQLProbe].
func (p *SQLProbe) SetName(name string) *SQLProbe {
	p.name = name

	return p
}

// Check returns a successful [healthcheck.CheckerProbeResult] if the database connection can be pinged.
func (p *SQLProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	err := p.db.PingContext(ctx)
	if err != nil {
		log.CtxLogger(ctx).Error().Err(err).Msg("database ping error")

		return healthcheck.NewCheckerProbeResult(false, fmt.Sprintf("database ping error: %v", err))
	}

	return healthcheck.NewCheckerProbeResult(true, "database ping success")
}
