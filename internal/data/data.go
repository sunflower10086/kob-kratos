package data

import (
	"database/sql"

	"kob-kratos/internal/conf"
	"kob-kratos/internal/data/gormgen/query"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewPostgresDB,
)

// Data .
type Data struct {
	DB *query.Query
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{
		DB: query.Use(db),
	}, cleanup, nil
}

func NewPostgresDB(c *conf.Data, logger log.Logger) (*gorm.DB, func(), error) {
	var (
		sqlDB *sql.DB
		db    *gorm.DB
		err   error
	)
	postgresConf := postgres.Config{
		DSN: c.GetDatabase().Dsn,
	}
	gormConfig := configLog(c.GetDatabase().LogMode, int(c.GetDatabase().CreateBatchSize))
	if db, err = gorm.Open(postgres.New(postgresConf), gormConfig); err != nil {
		return nil, nil, errors.Wrap(err, "opens database failed")
	}
	if sqlDB, err = db.DB(); err != nil {
		return nil, nil, errors.Wrap(err, "get database connection failed")
	}

	sqlDB.SetMaxIdleConns(int(c.GetDatabase().MaxIdleCons))
	sqlDB.SetMaxOpenConns(int(c.GetDatabase().MaxOpenCons))
	return db, func() {
		sqlDB.Close()
	}, nil
}

// configLog 根据配置决定是否开启日志
func configLog(mod bool, batchSize int) (c *gorm.Config) {
	c = &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不加复数形式，false默认加
		},
		CreateBatchSize: batchSize,
	}

	if mod {
		c.Logger = logger.Default.LogMode(logger.Info)
	}
	return
}
