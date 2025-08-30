package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gsmayya/theater/utils"
)

type Database struct {
	db *sql.DB
}

var (
	instance *Database
	once     sync.Once
)

// GetDatabase returns a singleton database instance with connection pooling
func GetDatabase() *Database {
	once.Do(func() {
		dbHost := utils.GetEnvOrDefault("DB_HOST", "localhost")
		dbUser := utils.GetEnvOrDefault("DB_USER", "user")
		dbPassword := utils.GetEnvOrDefault("DB_PASSWORD", "password")
		dbName := utils.GetEnvOrDefault("DB_NAME", "theater_booking")
		dbPort := utils.GetEnvOrDefault("DB_PORT", "3306")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// Configure connection pool
		db.SetMaxOpenConns(25)                 // Maximum number of open connections
		db.SetMaxIdleConns(25)                 // Maximum number of idle connections
		db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

		// Test the connection
		if err := db.Ping(); err != nil {
			log.Fatal("Failed to ping database:", err)
		}

		instance = &Database{db: db}
		log.Println("Database connection established with connection pooling")
	})

	return instance
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB returns the underlying sql.DB instance for complex queries
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// Ping tests the database connection
func (d *Database) Ping() error {
	return d.db.Ping()
}

// BeginTransaction starts a new database transaction
func (d *Database) BeginTransaction() (*sql.Tx, error) {
	return d.db.Begin()
}

// ExecuteInTransaction executes a function within a database transaction
func (d *Database) ExecuteInTransaction(fn func(*sql.Tx) error) error {
	tx, err := d.BeginTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
