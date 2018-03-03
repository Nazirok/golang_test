package store

import (
	"github.com/golang_test/handler"
)

type DbService interface {
	Set(key int, value interface{})
	Delete(key int) bool
	Get(key int) (handler.ClientBody, bool)
	GetAllDataJson() ([]byte, error)
}

type DbIdService interface {
	GenerateUnicId() interface{}
}
