package env

import "github.com/kelseyhightower/envconfig"

// Config holds environment paramenters.
type Config struct {
	Debug        bool   `envconfig:"debug" default:"false" desc:"Application Debug Mode"`
	Stage        string `envconfig:"stage" default:"development" desc:"Application Stage Environment"`
	Port         string `envconfig:"port" default:"5000" desc:"Application Port Number"`
	DatabaseURL  string `envconfig:"database_url" required:"true" desc:"Database URL"`
	PassesBucket string `envconfig:"passes_bucket" required:"true" desc:"Passes Bucket Name"`
	ServerSecret string `envconfig:"server_secret" required:"true" desc:"JWT Server Secret"`
}

// NewConfig parses and returns a new config.
func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return c, err
	}
	return c, nil
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
