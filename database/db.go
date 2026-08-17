package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db := stdlib.OpenDB(*config)

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database.")

	return db
}
