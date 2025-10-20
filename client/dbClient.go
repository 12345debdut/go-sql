package client

import (
	"fmt"
	"log"

	"go-sql/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbClient struct {
	providers map[string]*gorm.DB
}

func (db *DbClient) Connect(config models.SqlClientConnectionConfig) (*gorm.DB, error) {
	sslModeDb, err := models.SslModeDBVersion(config.SslMode)
	if err == nil {
		log.Fatal(err)
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Host, config.User, config.Password, config.DbName, config.Port, sslModeDb)
	ormDb, err := gorm.Open(postgres.Open(dsn), config.OrmConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = ormDb.AutoMigrate(config.Models...)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	db.providers[config.DbName] = ormDb
	return ormDb, nil
}

func (db *DbClient) Disconnect() error {
	for dbName, ormDb := range db.providers {
		log.Printf("Disconnecting from database: %s\n", dbName)
		sqlDb, err := ormDb.DB()
		if err != nil {
			log.Fatalf("Failed to disconnect from database: %v", err)
			return err
		}
		closeErr := sqlDb.Close()
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}
