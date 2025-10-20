package client

import (
	"go-sql/models"

	"gorm.io/gorm"
)

func NewSqlClient() SqlClient {
	return &DbClient{providers: make(map[string]*gorm.DB)}
}

type SqlClient interface {
	Connect(config models.SqlClientConnectionConfig) (*gorm.DB, error)
	Disconnect() error
}
