package client

import (
	"github.com/12345debdut/go-sql/models"

	"sync"

	"gorm.io/gorm"
)

func NewSqlClient() SqlClient {
	return &DbClient{}
}

type DbProvider struct {
	dbDriver *gorm.DB
	mutex    *sync.RWMutex
	config   models.SqlClientConnectionConfig
}

type SqlClient interface {
	Connect(config models.SqlClientConnectionConfig) (*gorm.DB, error)
	Disconnect() error
	ReadLock(config models.SqlClientConnectionConfig)
	ReadUnlock(config models.SqlClientConnectionConfig)
	WriteLock(config models.SqlClientConnectionConfig)
	WriteUnlock(config models.SqlClientConnectionConfig)
}
