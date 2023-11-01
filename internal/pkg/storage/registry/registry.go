package registry

import (
	"dps/internal/pkg/config"
	"dps/internal/pkg/storage"
	"dps/internal/pkg/storage/mysql"
	"dps/logger"
)

var conf = &config.Config

// StorageFactory is factory to get storage instance.
type StorageFactory interface {
	// GetStorage will return the Storage instance.
	GetStorage() storage.Store
}

var sqlFac = &SingletonFactory{
	creatorFunction: func() storage.Store {
		return &mysql.Store{}
	},
}

var storeFactorys = map[string]StorageFactory{
	"mysql":    sqlFac,
	"postgres": sqlFac,
}

// GetStore returns storage.Store
func GetStore() storage.Store {
	return storeFactorys[conf.Store.Driver].GetStorage()
}

// MustWaitStoreUp wait for db to go up
func MustWaitStoreUp() {
	for err := GetStore().Ping(); err != nil; {
		logger.FatalIfError(err)
	}
}
