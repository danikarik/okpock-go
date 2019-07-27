package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	tempPassesTable        = "CREATE TEMPORARY TABLE `passes` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `serial_number` VARCHAR(255) NOT NULL, `authentication_token` VARCHAR(255) NOT NULL, `pass_type_id` VARCHAR(255) NOT NULL, `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `serial_number_unique` (`serial_number`) );"
	tempRegistrationsTable = "CREATE TEMPORARY TABLE `registrations` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `uuid` VARCHAR(255) NOT NULL, `device_id` VARCHAR(255) NOT NULL, `push_token` VARCHAR(255) NOT NULL, `serial_number` VARCHAR(255) NOT NULL, `pass_type_id` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `serial_number_unique` (`serial_number`) );"
	tempLogsTable          = "CREATE TEMPORARY TABLE `logs` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `uuid` VARCHAR(191) NOT NULL, `remote_address` VARCHAR(191) NULL, `request_id` VARCHAR(191) NOT NULL, `message` VARCHAR(191) NOT NULL, `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `uuid_unique` (`uuid`) );"
	tempUsersTable         = "CREATE TEMPORARY TABLE `users` ( `id` INT(10) unsigned NOT NULL AUTO_INCREMENT, `role` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT \"client\", `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `email` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `password_hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `confirmed_at` TIMESTAMP NULL DEFAULT NULL, `invited_at` TIMESTAMP NULL DEFAULT NULL, `confirmation_token` VARCHAR(191) DEFAULT NULL, `confirmation_sent_at` TIMESTAMP NULL DEFAULT NULL, `recovery_token` VARCHAR(191) DEFAULT NULL, `recovery_sent_at` TIMESTAMP NULL DEFAULT NULL, `email_change_token` VARCHAR(191) DEFAULT NULL, `email_change` VARCHAR(191) DEFAULT NULL, `email_change_sent_at` TIMESTAMP NULL DEFAULT NULL, `last_signin_at` TIMESTAMP NULL DEFAULT NULL, `raw_app_metadata` text DEFAULT NULL, `raw_user_metadata` text DEFAULT NULL, `is_super_admin` TINYINT(1) DEFAULT 0, `created_at` TIMESTAMP NULL DEFAULT NULL, `updated_at` TIMESTAMP NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `users_username_unique` (`username`), UNIQUE KEY `users_email_unique` (`email`) );"
)

const (
	insertPassesTable        = "INSERT INTO `passes` (`serial_number`, `authentication_token`, `pass_type_id`, `updated_at`) VALUES ('%s', '%s', '%s', '%s');"
	insertRegistrationsTable = "INSERT INTO `registrations` (`uuid`, `device_id`, `push_token`, `serial_number`, `pass_type_id`) VALUES ('%s', '%s', '%s', '%s', '%s');"
	insertUsersTable         = "INSERT INTO `users` (`email`, `username`, `password_hash`) VALUES ('%s', '%s', '%s');"
)

func executeTempScripts(ctx context.Context, t *testing.T, schema, data []string) (*sqlx.DB, error) {
	env, err := env.NewLookup("TEST_DATABASE_URL")
	if err != nil {
		t.Skip(err)
	}

	conn, err := sqlx.Connect("mysql", env.Get("TEST_DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	for _, tab := range schema {
		_, err = conn.ExecContext(ctx, tab)
		if err != nil {
			return nil, err
		}
	}

	for _, sql := range data {
		_, err = conn.ExecContext(ctx, sql)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}
