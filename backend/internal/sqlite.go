package sqlite

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func CreateAllTables() *sql.DB {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	db.SetMaxOpenConns(1)

	log.Println("Running database migrations...")

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("failed to create sqlite driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations/sqlite",
		"sqlite3", driver,
	)
	if err != nil {
		log.Fatal("failed to initialize migrate:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if ver, dirty, vErr := m.Version(); vErr == nil && dirty {
			log.Printf("dirty at version %d, forcing back to version %d", ver, ver-1)
			if fErr := m.Force(int(ver - 1)); fErr != nil {
				log.Fatalf("force failed: %v", fErr)
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("migration retry failed: %v", err)
			}
		} else {
			log.Fatalf("migration failed: %v", err)
		}
	}

	log.Println("Database migrations completed successfully.")
	return db
}
