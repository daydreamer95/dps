package pkg

import (
	"dps/internal/pkg/storage"
	"dps/internal/pkg/storage/registry"
)

func GetStore() storage.Store {
	return registry.GetStore()
}
