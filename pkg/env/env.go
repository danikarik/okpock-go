package env

import (
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/apns"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/danikarik/okpock/pkg/mail"
	"github.com/danikarik/okpock/pkg/pkpass"
)

// New returns a new instance of `Env`.
func New(cfg Config, passkit api.PassKit, auth api.Auth, logic api.Logic,
	storage filestore.Storage, mailer mail.Mailer, coupon pkpass.Signer,
	notificator apns.Notificator) *Env {
	return &Env{
		Config:       cfg,
		PassKit:      passkit,
		Auth:         auth,
		Logic:        logic,
		Storage:      storage,
		Mailer:       mailer,
		CouponSigner: coupon,
		Notificator:  notificator,
	}
}

// Env holds stores and config.
type Env struct {
	Config       Config
	PassKit      api.PassKit
	Auth         api.Auth
	Logic        api.Logic
	Storage      filestore.Storage
	Mailer       mail.Mailer
	CouponSigner pkpass.Signer
	Notificator  apns.Notificator
}
