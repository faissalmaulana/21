package repository_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		log.Fatal("TEST_DB_DSN must be set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	testDB = db

	code := m.Run()

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		log.Printf("migration down failed: %v", err)
	}

	srcErr, dbErr := migrator.Close()
	if srcErr != nil {
		log.Printf("source close error: %v", srcErr)
	}
	if dbErr != nil {
		log.Printf("db close error: %v", dbErr)
	}

	os.Exit(code)
}

func seedProject(t *testing.T, db *sql.DB) string {
	t.Helper()

	var newProjectID string

	err := db.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", "New Project").Scan(&newProjectID)
	if err != nil {
		t.Fatal(err)
	}

	return newProjectID
}

// seedTask seeds without project_id inserted
func seedTask(t *testing.T, db *sql.DB) string {
	t.Helper()

	var newTaskID string

	err := db.QueryRow("INSERT INTO tasks(name,start_at) VALUES($1,$2) RETURNING id", "New Task", time.Now()).Scan(&newTaskID)
	if err != nil {
		t.Fatal(err)
	}

	return newTaskID
}
