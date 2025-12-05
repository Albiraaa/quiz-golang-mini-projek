package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var DB *sql.DB

func InitDB() {
	// PRIORITY: use DATABASE_URL if exists (Railway ready)
	dsn := os.Getenv("DATABASE_URL")

	// fallback ke DB manual mode
	if dsn == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		sslMode := os.Getenv("DB_SSLMODE")

		if dbPort == "" {
			dbPort = "5432"
		}
		if sslMode == "" {
			sslMode = "disable"
		}

		dsn = "postgres://" + dbUser + ":" + dbPass +
			"@" + dbHost + ":" + dbPort +
			"/" + dbName + "?sslmode=" + sslMode
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed open DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Failed ping DB:", err)
	}

	log.Println("🟢 DB connected!")

	autoMigrate()
}

func autoMigrate() {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations", // folder migrations wajib ada
	}

	// Always run UP migrations
	n, err := migrate.Exec(DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	if n > 0 {
		log.Printf("🚀 %d new migrations applied!\n", n)
	} else {
		log.Println("👌 No new migrations, already up to date.")
	}
}
