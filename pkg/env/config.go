package env

import "github.com/kelseyhightower/envconfig"

// NewConfig parses and returns a new config.
func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return c, err
	}
	return c, nil
}

// Config holds environment variables.
type Config struct {
	Debug        bool              `envconfig:"debug" default:"false" desc:"Application Debug Mode"`
	Stage        string            `envconfig:"stage" default:"development" desc:"Application Stage Environment"`
	Port         string            `envconfig:"port" default:"5000" desc:"Application Port Number"`
	DatabaseURL  string            `envconfig:"database_url" required:"true" desc:"Database URL"`
	UploadBucket string            `envconfig:"upload_bucket" required:"true" desc:"User Uploads Bucket Name"`
	PassesBucket string            `envconfig:"passes_bucket" required:"true" desc:"Passes Bucket Name"`
	ServerSecret string            `envconfig:"server_secret" required:"true" desc:"JWT Server Secret"`
	MailerRegion string            `envconfig:"mailer_region" required:"true" desc:"Mailer Region"`
	Certificates CertificateConfig `envconfig:"certificates" required:"true" desc:"Certificates Config"`
}

// CertificateConfig holds environment variables related to certificates.
type CertificateConfig struct {
	Team     string      `envconfig:"team" required:"true" desc:"Apple Team Identifier"`
	Bucket   string      `envconfig:"bucket" required:"true" desc:"Apple Certificates Bucket Name"`
	RootCert string      `envconfig:"root_cert" required:"true" desc:"Apple WWDR Certificate"`
	APS      Certificate `envconfig:"aps" required:"true" desc:"Apple Push Service Certificate"`
	Coupon   Certificate `envconfig:"coupon" required:"true" desc:"Coupon Certificate"`
}

// Certificate holds certificate path and password.
type Certificate struct {
	Path string `envconfig:"path" required:"true" desc:"Certificate Path"`
	Pass string `envconfig:"pass" required:"true" desc:"Certificate Password"`
}

// Usage returns usage for config instance.
func Usage(c Config) error {
	return envconfig.Usage("", &c)
}

// Addr returns http address to listen.
func (c Config) Addr() string {
	return ":" + c.Port
}

// IsDevelopment checks stage environment on `development` mode.
func (c Config) IsDevelopment() bool {
	return c.Stage == "development"
}

// IsProduction checks stage environment on `production` mode.
func (c Config) IsProduction() bool {
	return c.Stage == "production"
}
