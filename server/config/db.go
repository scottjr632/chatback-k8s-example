package config

import "os"

// DB holds all the database configurations
type DB struct {
	DBName   string
	Host     string
	Password string
	User     string
}

func newDB() *DB {
	db := &DB{}
	host := os.Getenv("DB_HOST")
	if host != "" {
		db.Host = host
	} else {
		db.Host = "localhost"
	}

	password := os.Getenv("DB_PASSWORD")
	if password != "" {
		db.Password = password
	} else {
		db.Password = "password"
	}

	user := os.Getenv("DB_USER")
	if user != "" {
		db.User = user
	} else {
		db.User = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName != "" {
		db.DBName = dbName
	} else {
		db.DBName = "postgres"
	}
	return db
}
