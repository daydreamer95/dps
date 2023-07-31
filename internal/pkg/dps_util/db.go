package dps_util

import (
	"database/sql"
	"dps/internal/pkg/config"
	"dps/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

// ModelBase model base for gorm to provide base fields
type ModelBase struct {
	CreateTime *time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime *time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

func getGormDialetor(driver string, dsn string) gorm.Dialector {
	if driver == DBTypePostgres {
		return postgres.Open(dsn)
	}
	if driver != DBTypeMysql {
		logger.Fatal(fmt.Sprintf("unknown driver: %s", driver))
	}
	return mysql.Open(dsn)
}

var dbs sync.Map

// DB provide more func over gorm.DB
type DB struct {
	*gorm.DB
}

// Must set must flag, panic when error occur
func (m *DB) Must() *DB {
	db := m.InstanceSet("ivy.must", true)
	return &DB{DB: db}
}

// ToSQLDB get the sql.DB
func (m *DB) ToSQLDB() *sql.DB {
	d, err := m.DB.DB()
	logger.Error(fmt.Sprintf("Error: [%s]", err))
	return d
}

type tracePlugin struct{}

func (op *tracePlugin) Name() string {
	return "tracePlugin"
}

func (op *tracePlugin) Initialize(db *gorm.DB) (err error) {
	before := func(db *gorm.DB) {
		db.InstanceSet("ivy.startTime", time.Now())
	}

	after := func(db *gorm.DB) {
		_ts, _ := db.InstanceGet("ivy.startTime")
		sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		logger.Debug(fmt.Sprintf("used: %d ms affected: %d sql is: %s", time.Since(_ts.(time.Time)).Milliseconds(), db.RowsAffected, sql))
		if v, ok := db.InstanceGet("ivy.must"); ok && v.(bool) {
			if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
				panic(db.Error)
			}
		}
	}

	beforeName := "cb_before"
	afterName := "cb_after"

	logger.Debug(fmt.Sprintf("installing db plugin: %s", op.Name()))
	// before
	_ = db.Callback().Create().Before("gorm:before_create").Register(beforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(beforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(beforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(beforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(beforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(beforeName, before)

	// after
	_ = db.Callback().Create().After("gorm:after_create").Register(afterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(afterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(afterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(afterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(afterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(afterName, after)
	return
}

// DbGet get db connection for specified conf
func DbGet(conf config.Store, ops ...func(*gorm.DB)) *DB {
	dsn := GetDsn(conf)
	db, ok := dbs.Load(dsn)
	if !ok {
		logger.Info(fmt.Sprintf("connecting '%s' '%s' '%s' '%d' '%s'", conf.Driver, conf.Host, conf.User, conf.Port, conf.Db))
		db1, err := gorm.Open(getGormDialetor(conf.Driver, dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		logger.FatalIfError(err)
		err = db1.Use(&tracePlugin{})
		logger.FatalIfError(err)
		db = &DB{DB: db1}
		for _, op := range ops {
			op(db1)
		}
		dbs.Store(dsn, db)
	}
	return db.(*DB)
}

// GetDsn get dsn from map config
func GetDsn(conf config.Store) string {
	host := conf.Host
	driver := conf.Driver
	dsn := map[string]string{
		"mysql": fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true",
			conf.User, conf.Password, host, conf.Port, conf.Db),
		"postgres": fmt.Sprintf("host=%s user=%s password=%s dbname='%s' search_path=%s port=%d sslmode=disable",
			host, conf.User, conf.Password, conf.Db, conf.Schema, conf.Port),
	}[driver]
	if dsn == "" {
		logger.Fatal(fmt.Sprintf("unknow driver: %s", driver))
	}
	return dsn
}
