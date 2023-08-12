package pkg

import (
	"dps/internal/pkg/storage"
	"dps/internal/pkg/storage/registry"
	"github.com/go-playground/validator/v10"
)

func init() {
	validate = validator.New()
}

var validate *validator.Validate

func GetStore() storage.Store {
	return registry.GetStore()
}
