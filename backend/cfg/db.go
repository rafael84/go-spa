package cfg

import (
	"fmt"
	"os"
)

type db struct {
	User     string
	Name     string
	Password string
	Host     string
	Port     string
	SSLmode  string
}

var DB db

func init() {
	DB.User = os.Getenv("DB_USER")
	DB.Name = os.Getenv("DB_NAME")
	DB.Password = os.Getenv("DB_PASSWORD")
	DB.Host = os.Getenv("DB_HOST")
	DB.Port = os.Getenv("DB_PORT")
	DB.SSLmode = os.Getenv("DB_SSLMODE")
}

func (db *db) ConnectionURL() string {
	format := `postgres://%s:%s@%s:%s/%s?sslmode=%s`
	return fmt.Sprintf(format, db.User, db.Password, db.Host, db.Port, db.Name, db.SSLmode)
}
