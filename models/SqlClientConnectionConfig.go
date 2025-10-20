package models

import (
	"log"

	"gorm.io/gorm"
)

type SqlClientConnectionConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	DbName    string
	SslMode   SslMode
	OrmConfig *gorm.Config
	Models    []interface{}
}

type SslMode int

const (
	Disable SslMode = iota
	Allow
)

type SslModeNotFoundError struct {
	error
}

func SslModeDBVersion(sslMode SslMode) (string, error) {
	switch sslMode {
	case Disable:
		return "disable", nil
	case Allow:
		return "allow", nil
	default:
		log.Fatal("Invalid sslMode")
		return "", SslModeNotFoundError{}
	}
}
