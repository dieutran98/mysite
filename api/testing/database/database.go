package database

import (
	"context"
	"database/sql"
	"fmt"
	"mysite/constants"
	"mysite/migration"
	"mysite/pkgs/env"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	user         string = "testUser"
	password     string = "testPassword"
	host         string = "localhost"
	port         string = "1234"
	databaseName string = "testDatabase"
)

func SetupDatabaseForTesting() (*dockertest.Pool, *dockertest.Resource, error) {
	pool, err := createDockerPool()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create docker pool")
	}

	resource, err := createDockerResource(pool)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create resource")
	}

	defer func() { // remove resource if failed create
		if err != nil {
			if err := PurgeResource(pool, resource); err != nil {
				fmt.Println("failed to purge resource")
			}
		}
	}()

	if err := updateEnv(resource); err != nil {
		return nil, nil, errors.Wrap(err, "failed to update env resource")
	}

	if err = checkIsDatabaseRunning(pool); err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup database")
	}

	if err := migration.Migrate(connectUrl()); err != nil {
		return nil, nil, errors.Wrap(err, "failed to migrate")
	}

	return pool, resource, nil
}

func PurgeResource(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if pool == nil {
		return errors.New("pool is empty")
	}
	if resource == nil {
		return errors.New("resource is empty")
	}

	return pool.Purge(resource)
}

func createDockerPool() (*dockertest.Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup docker test pool")
	}

	if err := pool.Client.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed connect to docker")
	}

	return pool, nil
}

func createDockerResource(pool *dockertest.Pool) (*dockertest.Resource, error) {
	if pool == nil {
		return nil, errors.New("pool is empty")
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.0-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", user),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			fmt.Sprintf("POSTGRES_DB=%s", databaseName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to run container")
	}

	return resource, nil
}

func checkIsDatabaseRunning(pool *dockertest.Pool) error {
	if pool == nil {
		return errors.New("pool is empty")
	}

	return pool.Retry(func() error {
		if _, err := connectDb(); err != nil {
			return err
		}
		return nil
	})
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

func updateEnv(resource *dockertest.Resource) error {
	if err := env.ReadEnv(func(appEnv *env.AppEnv) {
		appEnv.Database.User = user
		appEnv.Database.Password = password
		appEnv.Database.Database = databaseName
		appEnv.Database.SslMode = "disable"
		appEnv.Database.ConnMaxLifeIdle = 30
		appEnv.Database.ConnMaxLifeTime = 30
		appEnv.Database.ConnMaxOpen = 20
		appEnv.Database.TransactionTimeout = 30
		appEnv.Database.HostName = resource.GetBoundIP("5432/tcp")
		appEnv.Database.Port = resource.GetPort("5432/tcp")
	}); err != nil {
		return errors.Wrap(err, "failed to read env")
	}

	return nil
}

func SetTestTransactionCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, constants.Testing, true)
	return ctx
}
