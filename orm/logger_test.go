package orm_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/orm"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	base "gorm.io/gorm"
	baselogger "gorm.io/gorm/logger"
)

type model struct {
	Name string
}

func TestLoggerError(t *testing.T) {
	t.Parallel()

	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.ErrorLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	ctx := logger.WithContext(context.Background())

	orm.NewCtxOrmLogger(baselogger.Error, false).Error(ctx, "message")

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "error",
		"message": "message",
	})
}

func TestLoggerWarn(t *testing.T) {
	t.Parallel()

	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.WarnLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	ctx := logger.WithContext(context.Background())

	orm.NewCtxOrmLogger(baselogger.Warn, false).Warn(ctx, "message")

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "warn",
		"message": "message",
	})
}

func TestLoggerInfo(t *testing.T) {
	t.Parallel()

	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.InfoLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	ctx := logger.WithContext(context.Background())

	orm.NewCtxOrmLogger(baselogger.Info, false).Info(ctx, "message")

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "message",
	})
}

func TestLoggerTraceWithValues(t *testing.T) {
	t.Parallel()

	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.TraceLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDsn(":memory:"),
		orm.WithDriver(orm.Sqlite),
		orm.WithConfig(base.Config{
			Logger: orm.NewCtxOrmLogger(baselogger.Info, true),
		}),
	)
	assert.NoError(t, err)

	db = db.WithContext(logger.WithContext(context.Background()))
	err = db.AutoMigrate(&model{})
	assert.NoError(t, err)

	cases := []struct {
		run  func() error
		sql  string
		fail bool
	}{
		{
			run: func() error {
				return db.Create(&model{Name: "test"}).Error
			},
			sql:  "INSERT INTO `models` (`name`) VALUES (\"test\")",
			fail: false,
		},
		{
			run: func() error {
				return db.Model(&model{}).Find(&[]*model{}).Error
			},
			sql:  "SELECT * FROM `models`",
			fail: false,
		},
		{
			run: func() error {
				return db.Where(&model{Name: "test"}).First(&model{}).Error
			},
			sql:  "SELECT * FROM `models` WHERE `models`.`name` = \"test\" ORDER BY `models`.`name` LIMIT 1",
			fail: false,
		},
		{
			run: func() error {
				return db.Where(&model{Name: "invalid"}).First(&model{}).Error
			},
			sql:  "SELECT * FROM `models` WHERE `models`.`name` = \"invalid\" ORDER BY `models`.`name` LIMIT 1",
			fail: true,
		},
		{
			run: func() error {
				return db.Raw("invalid query").Scan(&model{}).Error
			},
			sql:  "invalid query",
			fail: true,
		},
	}

	for _, c := range cases {
		buffer.Reset()

		err = c.run()

		if err != nil && !c.fail {
			t.Errorf("unexpected error: %s (%T)", err, err)
		}

		level := "debug"
		if c.fail {
			level = "error"
		}

		logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
			"level":    level,
			"sqlQuery": c.sql,
		})
	}
}

func TestLoggerTraceWithoutValues(t *testing.T) {
	t.Parallel()

	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.TraceLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDsn(":memory:"),
		orm.WithDriver(orm.Sqlite),
		orm.WithConfig(base.Config{
			Logger: orm.NewCtxOrmLogger(baselogger.Info, false),
		}),
	)
	assert.NoError(t, err)

	db = db.WithContext(logger.WithContext(context.Background()))
	err = db.AutoMigrate(&model{})
	assert.NoError(t, err)

	err = db.Create(&model{Name: "test"}).Error
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":    "debug",
		"sqlQuery": "INSERT INTO `models` (`name`) VALUES (?)",
	})
}

func TestLoggerTraceSilent(t *testing.T) {
	t.Parallel()

	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.TraceLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	ormLogger := orm.NewCtxOrmLogger(baselogger.Info, false)
	ormLogger.LogMode(baselogger.Silent)

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDsn(":memory:"),
		orm.WithDriver(orm.Sqlite),
		orm.WithConfig(base.Config{
			Logger: ormLogger,
		}),
	)
	assert.NoError(t, err)

	db = db.WithContext(logger.WithContext(context.Background()))
	err = db.AutoMigrate(&model{})
	assert.NoError(t, err)

	err = db.Create(&model{Name: "test"}).Error
	assert.NoError(t, err)

	records, _ := buffer.Records()
	assert.Len(t, records, 0)
}
