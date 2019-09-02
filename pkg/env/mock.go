package env

import (
	"os"

	fsmock "github.com/danikarik/okpock/pkg/filestore/memory"
	mlmock "github.com/danikarik/okpock/pkg/mail/memory"
	"github.com/danikarik/okpock/pkg/store/sequel"
	_ "github.com/go-sql-driver/mysql" //
	"github.com/jmoiron/sqlx"
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

	conn, err := sqlx.Connect("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	cleanUp := []string{
		"DELETE FROM `user_uploads`",
		"DELETE FROM `uploads`",
		"DELETE FROM `user_projects`",
		"DELETE FROM `projects`",
		"DELETE FROM `users`",
		"DELETE FROM `logs`",
		"DELETE FROM `registrations`",
		"DELETE FROM `passes`",
	}

	for _, script := range cleanUp {
		_, err = conn.Exec(script)
		if err != nil {
			return nil, err
		}
	}

	db := sequel.New(conn)
	fs := fsmock.New()
	ml := mlmock.New()

	return New(cfg, db, db, db, fs, ml), nil
}
