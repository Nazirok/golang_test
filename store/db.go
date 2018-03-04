package store

import (
	"github.com/golang_test/handler"
)

type DbService interface {
	Set(key int, value handler.ClientBody)
	Delete(key int) bool
	Get(key int) (handler.ClientBody, bool)
	GetAllDataJson() ([]byte, error)
}
