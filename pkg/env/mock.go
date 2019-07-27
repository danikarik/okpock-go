package env

import (
	"os"

	fsmock "github.com/danikarik/okpock/pkg/filestore/memory"
	dbmock "github.com/danikarik/okpock/pkg/store/memory"
)

// NewMock returns a new mock `Env`.
func NewMock() *Env {
	cfg := Config{
		Stage:        "test",
		Port:         "5000",
		DatabaseURL:  os.Getenv("TEST_DATABASE_URL"),
		PassesBucket: os.Getenv("TEST_PASSES_BUCKET"),
		ServerSecret: os.Getenv("TEST_SERVER_SECRET"),
	}

	db := dbmock.New()
	fs := fsmock.New()

	return New(cfg, db, db, fs)
}
