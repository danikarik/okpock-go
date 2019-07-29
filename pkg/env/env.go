package env

import (
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/danikarik/okpock/pkg/mail"
)

// Env holds stores and config.
type Env struct {
	Config  Config
	PassKit api.PassKit
	Auth    api.Auth
	Storage filestore.Storage
	Mailer  mail.Mailer
}

// New returns a new instance of `Env`.
func New(c Config, passkit api.PassKit, auth api.Auth, storage filestore.Storage, mailer mail.Mailer) *Env {
	return &Env{
		Config:  c,
		PassKit: passkit,
		Auth:    auth,
		Storage: storage,
		Mailer:  mailer,
	}
}
