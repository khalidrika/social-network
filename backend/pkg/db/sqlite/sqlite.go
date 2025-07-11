package db

import (
	"database/sql"
	"fmt"
	"log"

	"social-network/pkg/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// initialize db
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./tempDB.db")
	if err != nil {
		log.Fatal(err)
	}

	err = Migrations(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// InitRepositories should be called in server.go
// creates connection to db for all rep
func InitRepositories(db *sql.DB) *models.Repositories {
	return &models.Repositories{
		UserRepo:    &UserRepository{DB: db},
		SessionRepo: &SessionRepository{DB: db},
		GroupRepo:   &GroupRepository{DB: db},
		PostRepo:    &PostRepository{DB: db},
		CommentRepo: &CommentRepository{DB: db},
		NotifRepo:   &NotifRepository{DB: db},
		EventRepo:   &EventRepository{DB: db},
		MsgRepo:     &MsgRepository{DB: db},
	}
}

func Migrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../backend/pkg/db/migration/sqlite",
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
