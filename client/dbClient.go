package client

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/12345debdut/go-sql/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DbClient Used sync map to store the database connection
// This should be concurrent safe so that multiple goroutines can access the same database concurrently
type DbClient struct {
	providers sync.Map
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
	db.providers.Store(config.DbName, &DbProvider{config: config, dbDriver: ormDb, mutex: &sync.RWMutex{}})
	return ormDb, nil
}

func (db *DbClient) Disconnect() error {
	var resultError error
	db.providers.Range(func(key, value interface{}) bool {
		dbName, dbNameOk := key.(string)
		if dbNameOk == false {
			resultError = errors.New("Key is not a string")
			return false
		}
		dbProvider, dbValueOk := value.(*DbProvider)
		if dbValueOk == false {
			resultError = errors.New("Value is not a DbProvider")
			return false
		}
		log.Printf("Disconnecting from database: %s\n", dbName)
		sqlDb, err := dbProvider.dbDriver.DB()
		if err != nil {
			resultError = err
			log.Fatalf("Failed to disconnect from database: %v", err)
			return false
		}
		closeErr := sqlDb.Close()
		if closeErr != nil {
			resultError = err
			return false
		}
		return true
	})
	if resultError != nil {
		return resultError
	}
	return nil
}

func (db *DbClient) ReadLock(config models.SqlClientConnectionConfig) {
	provider, ok := db.providers.Load(config.DbName)
	if ok == false {
		return
	}
	dbProvider, ok := provider.(*DbProvider)
	if ok == false {
		return
	}
	dbProvider.mutex.RLock()
}

func (db *DbClient) ReadUnlock(config models.SqlClientConnectionConfig) {
	provider, ok := db.providers.Load(config.DbName)
	if ok == false {
		return
	}
	dbProvider, ok := provider.(*DbProvider)
	if ok == false {
		return
	}
	dbProvider.mutex.RUnlock()
}

func (db *DbClient) WriteLock(config models.SqlClientConnectionConfig) {
	provider, ok := db.providers.Load(config.DbName)
	if ok == false {
		return
	}
	dbProvider, ok := provider.(*DbProvider)
	if ok == false {
		return
	}
	dbProvider.mutex.Lock()
}

func (db *DbClient) WriteUnlock(config models.SqlClientConnectionConfig) {
	provider, ok := db.providers.Load(config.DbName)
	if ok == false {
		return
	}
	dbProvider, ok := provider.(*DbProvider)
	if ok == false {
		return
	}
	dbProvider.mutex.Unlock()
}
