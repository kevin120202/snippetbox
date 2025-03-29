package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kevin120202/snippetbox/internal/models"
	_ "github.com/lib/pq"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		logger.Error("DB_URL must be set")
	}

	// Connect to postgres
	db, err := openDB(dbURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Connection pool close before main() exits
	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	// Initialize the pool
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create a connection pool and check for any errors
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
