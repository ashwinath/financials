package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required only for testing
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host             = getEnvOrDefault("DB_HOST", "localhost")
	user             = getEnvOrDefault("DB_USER", "postgres")
	password         = getEnvOrDefault("DB_PASSWORD", "postgres")
	database         = getEnvOrDefault("DB_NAME", "postgres")
	port             = getEnvOrDefault("DB_PORT", "5432")
	timeZone         = getEnvOrDefault("DB_TIMEZONE", "Asia/Singapore")
	migrationsFolder = "../migrations"
)

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value
}

func connectionString(db string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname='%s' port=%s sslmode=disable TimeZone=%s",
		host,
		user,
		password,
		db,
		port,
		timeZone,
	)
}

func create(conn *sql.DB, dbName string) (*sql.DB, error) {
	if _, err := conn.Exec("CREATE DATABASE " + dbName); err != nil {
		return nil, err
	} else if testDb, err := sql.Open("postgres", connectionString(dbName)); err != nil {
		if _, err := conn.Exec("DROP DATABASE " + dbName); err != nil {
			log.Fatalf("Failed to cleanup integration test database: \n%s", err)
		}
		return nil, err
	} else {
		return testDb, nil
	}
}

func migrate(db *sql.DB, dbName string) (*sql.DB, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	defer driver.Close()

	if migrations, err := gomigrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsFolder),
		dbName, driver); err != nil {
		return db, err
	} else if err = migrations.Up(); err != nil {
		return db, err
	}
	return sql.Open("postgres", connectionString(dbName))
}

func createTestDatabase() (*gorm.DB, func(), error) {
	testDbName := fmt.Sprintf("test_%d", time.Now().UnixNano())

	connStr := connectionString(database)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}

	testDb, err := create(conn, testDbName)
	if err != nil {
		return nil, nil, err
	}

	var gormDb *gorm.DB

	cleanup := func() {
		if gormDb != nil {
			db, err := gormDb.DB()
			if err != nil {
				log.Fatalf("failed to get db from gorm.")
			} else {
				db.Close()
			}
		}

		if err := testDb.Close(); err != nil {
			log.Fatalf("Failed to close connection to integration test database: \n%s", err)
		} else if _, err := conn.Exec("DROP DATABASE " + testDbName); err != nil {
			log.Fatalf("Failed to cleanup integration test database: \n%s", err)
		} else if err = conn.Close(); err != nil {
			log.Fatalf("Failed to close database: \n%s", err)
		}
	}

	if testDb, err = migrate(testDb, testDbName); err != nil {
		cleanup()
		return nil, nil, err
	} else if gormDb, err = gorm.Open(gormpg.Open(connectionString(testDbName)), &gorm.Config{}); err != nil {
		cleanup()
		return nil, nil, err
	} else {
		log.Printf(testDbName)
		return gormDb, cleanup, nil
	}
}

// WithTestDatabase is a method to test functions while creating a database that will be trashed later
func WithTestDatabase(t *testing.T, test func(t *testing.T, db *gorm.DB)) {
	if testDb, cleanupFn, err := createTestDatabase(); err != nil {
		t.Fatalf("Fail to create an integration test database: \n%s", err)
	} else {
		test(t, testDb)
		cleanupFn()
	}
}
