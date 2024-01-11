package healthcheck

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log"
	"gorm.io/gorm"
)

// DefaultProbeName is the name of the ORM probe.
const DefaultProbeName = "orm"

// OrmProbe is a probe compatible with the [healthcheck] module.
//
// [healthcheck]: https://github.com/ankorstore/yokai/tree/main/healthcheck
type OrmProbe struct {
	name string
	db   *gorm.DB
}

// NewOrmProbe returns a new [OrmProbe].
func NewOrmProbe(db *gorm.DB) *OrmProbe {
	return &OrmProbe{
		name: DefaultProbeName,
		db:   db,
	}
}

// NewOrmProbe returns the name of the [OrmProbe].
func (p *OrmProbe) Name() string {
	return p.name
}

// SetName sets the name of the [OrmProbe].
func (p *OrmProbe) SetName(name string) *OrmProbe {
	p.name = name

	return p
}

// Check returns a successful [healthcheck.CheckerProbeResult] if the database connection can be pinged.
func (p *OrmProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	db, err := p.db.DB()
	if err != nil {
		return healthcheck.NewCheckerProbeResult(false, fmt.Sprintf("database fetch error: %v", err))
	}

	err = db.Ping()
	if err != nil {
		log.CtxLogger(ctx).Error().Err(err).Msg("database ping error")

		return healthcheck.NewCheckerProbeResult(false, fmt.Sprintf("database ping error: %v", err))
	}

	return healthcheck.NewCheckerProbeResult(true, "database ping success")
}
