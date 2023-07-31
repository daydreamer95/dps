package mysql

import (
	"dps/internal/pkg/config"
	"dps/internal/pkg/dps_util"
	"gorm.io/gorm"
	"time"
)

var conf = &config.Config

// Store implements storage.Store, and storage with db
type Store struct {
}

func (s *Store) Ping() error {
	return dbGet().ToSQLDB().Ping()
}

// SetDBConn sets db conn pool
func SetDBConn(db *gorm.DB) {
	sqldb, _ := db.DB()
	sqldb.SetMaxOpenConns(int(conf.Store.MaxOpenConns))
	sqldb.SetMaxIdleConns(int(conf.Store.MaxIdleConns))
	sqldb.SetConnMaxLifetime(time.Duration(conf.Store.ConnMaxLifeTime) * time.Minute)
}

func dbGet() *dps_util.DB {
	return dps_util.DbGet(conf.Store, SetDBConn)
}
