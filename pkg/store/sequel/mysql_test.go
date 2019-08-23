package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	tempPassesTable        = "CREATE TABLE IF NOT EXISTS `passes` ( `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `authentication_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `passes_serial_number_unique` (`serial_number`) );"
	tempRegistrationsTable = "CREATE TABLE IF NOT EXISTS `registrations` ( `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `device_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `push_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `registrations_serial_number_unique` (`serial_number`) );"
	tempLogsTable          = "CREATE TABLE IF NOT EXISTS `logs` ( `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `remote_address` VARCHAR(191) COLLATE utf8mb4_unicode_ci NULL, `request_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `message` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (`id`) );"
	tempUsersTable         = "CREATE TABLE IF NOT EXISTS `users` ( `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `role` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT \"client\", `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `email` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `password_hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `confirmed_at` TIMESTAMP NULL DEFAULT NULL, `invited_at` TIMESTAMP NULL DEFAULT NULL, `confirmation_token` VARCHAR(191) DEFAULT \"\", `confirmation_sent_at` TIMESTAMP NULL DEFAULT NULL, `recovery_token` VARCHAR(191) DEFAULT \"\", `recovery_sent_at` TIMESTAMP NULL DEFAULT NULL, `email_change_token` VARCHAR(191) DEFAULT \"\", `email_change` VARCHAR(191) DEFAULT \"\", `email_change_sent_at` TIMESTAMP NULL DEFAULT NULL, `last_signin_at` TIMESTAMP NULL DEFAULT NULL, `raw_app_metadata` text DEFAULT NULL, `raw_user_metadata` text DEFAULT NULL, `is_super_admin` TINYINT(1) DEFAULT 0, `created_at` TIMESTAMP NULL DEFAULT NOW(), `updated_at` TIMESTAMP NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `users_username_unique` (`username`), UNIQUE KEY `users_email_unique` (`email`), UNIQUE KEY `users_username_and_email_unique` (`username`, `email`) );"
	tempOrganizationsTable = "CREATE TABLE IF NOT EXISTS `organizations` ( `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `user_id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `title` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `raw_metadata` text DEFAULT NULL, `created_at` TIMESTAMP NULL DEFAULT NOW(), `updated_at` TIMESTAMP NULL DEFAULT NOW(), PRIMARY KEY (`id`), FOREIGN KEY `organizations_user_reference` (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE, UNIQUE KEY `organizations_title_and_user_unique` (`user_id`, `title`) );"
	tempProjectsTable      = "CREATE TABLE IF NOT EXISTS `projects` ( `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `organization_id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL, `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `pass_type` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `background_image` VARCHAR(191) DEFAULT \"\", `footer_image` VARCHAR(191) DEFAULT \"\", `icon_image` VARCHAR(191) DEFAULT \"\", `strip_image` VARCHAR(191) DEFAULT \"\", `created_at` TIMESTAMP NULL DEFAULT NOW(), `updated_at` TIMESTAMP NULL DEFAULT NOW(), PRIMARY KEY (`id`), FOREIGN KEY `projects_organization_reference` (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE, UNIQUE KEY `projects_title_and_user_unique` (`organization_id`, `description`, `pass_type`) );"
)

const (
	insertPassesTable        = "INSERT INTO `passes` (`id`, `serial_number`, `authentication_token`, `pass_type_id`, `updated_at`) VALUES ('%s', '%s', '%s', '%s', '%s');"
	insertRegistrationsTable = "INSERT INTO `registrations` (`id`, `device_id`, `push_token`, `serial_number`, `pass_type_id`) VALUES ('%s', '%s', '%s', '%s', '%s');"
	insertUsersTable         = "INSERT INTO `users` (`id`, `email`, `username`, `password_hash`) VALUES ('%s', '%s', '%s', '%s');"
	insertOrganizationsTable = "INSERT INTO `organizations` (`id`, `user_id`, `title`, `description`) VALUES ('%s', '%s', '%s', '%s');"
	insertProjectsTable      = "INSERT INTO `projects` (`id`, `organization_id`, `description`, `pass_type`) VALUES ('%s', '%s', '%s', '%s');"
)

var cleanupTables = []string{
	"DROP TABLE IF EXISTS `projects`;",
	"DROP TABLE IF EXISTS `organizations`;",
	"DROP TABLE IF EXISTS `users`;",
	"DROP TABLE IF EXISTS `registrations`;",
	"DROP TABLE IF EXISTS `passes`;",
}

func executeTempScripts(ctx context.Context, t *testing.T, schema, data []string) (*sqlx.DB, error) {
	env, err := env.NewLookup("TEST_DATABASE_URL")
	if err != nil {
		t.Skip(err)
	}

	conn, err := sqlx.Connect("mysql", env.Get("TEST_DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	for _, tab := range cleanupTables {
		_, err = conn.ExecContext(ctx, tab)
		if err != nil {
			return nil, err
		}
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
