package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var DB *sql.DB

func InitDB() {
	// Contoh pakai env (sesuaikan sama Railway lu)
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE") // biasanya "disable" di lokal

	if dbPort == "" {
		dbPort = "5432"
	}
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}

	log.Println("DB connected")

	// Jalankan migration
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("failed to run migrations:", err)
	}
	log.Printf("Applied %d migrations\n", n)
}
