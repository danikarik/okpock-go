package env

import (
	"os"

	fsmock "github.com/danikarik/okpock/pkg/filestore/memory"
	mlmock "github.com/danikarik/okpock/pkg/mail/memory"
	dbmock "github.com/danikarik/okpock/pkg/store/memory"
)

// NewMock returns a new mock `Env`.
func NewMock() (*Env, error) {
	cfg := Config{
		Stage:        "test",
		Port:         "5000",
		DatabaseURL:  os.Getenv("TEST_DATABASE_URL"),
		UploadBucket: os.Getenv("TEST_UPLOAD_BUCKET"),
		PassesBucket: os.Getenv("TEST_PASSES_BUCKET"),
		ServerSecret: os.Getenv("TEST_SERVER_SECRET"),
		MailerRegion: os.Getenv("TEST_MAILER_REGION"),
	}

	db := dbmock.New()
	fs := fsmock.New()
	ml := mlmock.New()

	// TODO: add mock signer
	return New(cfg, db, db, db, fs, ml, nil), nil
}
