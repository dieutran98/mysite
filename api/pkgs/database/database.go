package database

import (
	"database/sql"
	"fmt"
	"mysite/pkgs/env"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

type DB struct {
	db *sql.DB
}

var (
	internalDB DB
	once       sync.Once
)

func SetupDatabase() error {
	var err error
	once.Do(func() {
		err = newBoilerDb()
	})
	return err
}

func Close() error {
	return internalDB.db.Close()
}

func getDb() *DB {
	if internalDB.db == nil {
		db, err := connectDb()
		if err != nil {
			panic(err)
		}

		internalDB.db = db
	}
	return &internalDB
}

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("pgx", connectUrl())
	if err != nil {
		return nil, errors.Wrap(err, "failed open database")
	}

	dbEnv := env.GetEnv().Database

	db.SetConnMaxIdleTime(time.Duration(dbEnv.ConnMaxLifeIdle) * time.Second)
	db.SetConnMaxLifetime(time.Duration(dbEnv.ConnMaxLifeTime) * time.Second)
	db.SetMaxIdleConns(dbEnv.ConnMaxOpen)
	db.SetMaxOpenConns(dbEnv.ConnMaxOpen)

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed ping database")
	}

	return db, nil
}

func connectUrl() string {
	dbEnv := env.GetEnv().Database
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbEnv.User, dbEnv.Password, dbEnv.HostName, dbEnv.Port, dbEnv.Database, dbEnv.SslMode)
}
