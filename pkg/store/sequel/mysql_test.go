package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

var clean = []string{
	"DELETE FROM `project_pass_cards`",
	"DELETE FROM `pass_cards`",
	"DELETE FROM `user_uploads`",
	"DELETE FROM `uploads`",
	"DELETE FROM `user_projects`",
	"DELETE FROM `projects`",
	"DELETE FROM `users`",
	"DELETE FROM `logs`",
	"DELETE FROM `registrations`",
	"DELETE FROM `passes`",
}

func testConnection(ctx context.Context, t *testing.T) (*sqlx.DB, error) {
	env, err := env.NewLookup("TEST_DATABASE_URL")
	if err != nil {
		t.Skip(err)
	}

	conn, err := sqlx.Connect("mysql", env.Get("TEST_DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	for _, script := range clean {
		_, err = conn.Exec(script)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func fakeUsername() string {
	return uuid.NewV4().String()
}

func fakeEmail() string {
	return uuid.NewV4().String() + "@example.com"
}

func fakeString() string {
	return uuid.NewV4().String()
}
