package env

import (
	"io/ioutil"
	"os"

	"github.com/danikarik/okpock/pkg/apns"
	fsmock "github.com/danikarik/okpock/pkg/filestore/memory"
	mlmock "github.com/danikarik/okpock/pkg/mail/memory"
	"github.com/danikarik/okpock/pkg/pkpass"
	dbmock "github.com/danikarik/okpock/pkg/store/memory"
	uuid "github.com/satori/go.uuid"
)

// NewMock returns a new mock `Env`.
func NewMock() (*Env, error) {
	cfg := Config{
		Stage:        "development",
		Port:         "5000",
		DatabaseURL:  os.Getenv("TEST_DATABASE_URL"),
		UploadBucket: os.Getenv("TEST_UPLOAD_BUCKET"),
		PassesBucket: os.Getenv("TEST_PASSES_BUCKET"),
		ServerSecret: os.Getenv("TEST_SERVER_SECRET"),
		MailerRegion: os.Getenv("TEST_MAILER_REGION"),
		Certificates: CertificateConfig{Team: uuid.NewV4().String()},
	}

	db := dbmock.New()
	fs := fsmock.New()
	ml := mlmock.New()

	var (
		rootCertPath   = os.Getenv("TEST_CERTIFICATES_ROOT_CERT")
		couponCertPath = os.Getenv("TEST_CERTIFICATES_COUPON_PATH")
		couponCertPass = os.Getenv("TEST_CERTIFICATES_COUPON_PASS")
	)

	rootCert, err := ioutil.ReadFile(rootCertPath)
	if err != nil {
		return nil, err
	}

	couponCert, err := ioutil.ReadFile(couponCertPath)
	if err != nil {
		return nil, err
	}

	couponSigner, err := pkpass.NewSigner(rootCert, couponCert, couponCertPass)
	if err != nil {
		return nil, err
	}

	return New(cfg, db, db, db, fs, ml, couponSigner, apns.NewMock()), nil
}
