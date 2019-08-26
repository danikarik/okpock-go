package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/danikarik/okpock/pkg/filestore/awsstore"
	"github.com/danikarik/okpock/pkg/mail"
	"github.com/danikarik/okpock/pkg/mail/awsmail"
	"github.com/danikarik/okpock/pkg/service"
	"github.com/danikarik/okpock/pkg/store/redistore"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

// Version holds build version.
var Version string

func main() {
	var err error

	var cfg env.Config
	{
		cfg, err = env.NewConfig()
		if err != nil {
			errorExit("%v", env.Usage(cfg))
		}
	}

	var logger *zap.Logger
	{
		var logCfg zap.Config
		if cfg.IsDevelopment() {
			logCfg = zap.NewDevelopmentConfig()
		} else {
			logCfg = zap.NewProductionConfig()
		}

		logCfg.DisableCaller = true
		logCfg.DisableStacktrace = true

		logger, err = logCfg.Build()
		if err != nil {
			errorExit("zap logger: %v", err)
		}
	}
	defer logger.Sync()

	var pool *redis.Pool
	{
		dialOptions := []redis.DialOption{
			redis.DialPassword(cfg.RedisPass),
			redis.DialConnectTimeout(15 * time.Second),
			redis.DialReadTimeout(15 * time.Second),
			redis.DialWriteTimeout(15 * time.Second),
		}
		pool = &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", cfg.RedisHost, dialOptions...)
			},
			MaxIdle:   5,
			MaxActive: 100,
			Wait:      true,
		}
	}
	defer pool.Close()

	var s3 filestore.Storage
	{
		s3, err = awsstore.New()
		if err != nil {
			errorExit("aws store: %v", err)
		}
	}

	var mailer mail.Mailer
	{
		mailer, err = awsmail.New(cfg.MailerRegion)
		if err != nil {
			errorExit("aws mail: %v", err)
		}
	}

	var srv *service.Service
	{
		db := redistore.New(pool)
		env := env.New(cfg, db, db, db, s3, mailer)

		srv = service.New(Version, env, logger)
	}

	logger.Info("server", zap.String("http_address", cfg.Addr()))
	errorExit("server: %v", http.ListenAndServe(cfg.Addr(), srv))
}

func errorExit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
