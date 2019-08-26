package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

const (
	tempPassesTable        = "CREATE TEMPORARY TABLE `passes` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `authentication_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `passes_serial_number_unique_idx` (`serial_number`) );"
	tempRegistrationsTable = "CREATE TEMPORARY TABLE `registrations` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `device_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `push_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `registrations_serial_number_unique_idx` (`serial_number`) );"
	tempLogsTable          = "CREATE TEMPORARY TABLE `logs` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `remote_address` VARCHAR(191) COLLATE utf8mb4_unicode_ci NULL, `request_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `message` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (`id`) );"
	tempUsersTable         = "CREATE TEMPORARY TABLE `users` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `role` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT \"client\", `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `email` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `password_hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `confirmed_at` TIMESTAMP NULL DEFAULT NULL, `invited_at` TIMESTAMP NULL DEFAULT NULL, `confirmation_token` VARCHAR(191) DEFAULT \"\", `confirmation_sent_at` TIMESTAMP NULL DEFAULT NULL, `recovery_token` VARCHAR(191) DEFAULT \"\", `recovery_sent_at` TIMESTAMP NULL DEFAULT NULL, `email_change_token` VARCHAR(191) DEFAULT \"\", `email_change` VARCHAR(191) DEFAULT \"\", `email_change_sent_at` TIMESTAMP NULL DEFAULT NULL, `last_signin_at` TIMESTAMP NULL DEFAULT NULL, `raw_app_metadata` text DEFAULT NULL, `raw_user_metadata` text DEFAULT NULL, `is_super_admin` TINYINT(1) DEFAULT 0, `created_at` TIMESTAMP NULL DEFAULT NOW(), `updated_at` TIMESTAMP NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `users_username_unique_idx` (`username`), UNIQUE KEY `users_email_unique_idx` (`email`), UNIQUE KEY `users_username_and_email_unique_idx` (`username`, `email`) );"
	tempProjectsTable      = "CREATE TEMPORARY TABLE `projects` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `title` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `organization_name` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `pass_type` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL, `background_image` VARCHAR(191) DEFAULT \"\", `footer_image` VARCHAR(191) DEFAULT \"\", `icon_image` VARCHAR(191) DEFAULT \"\", `strip_image` VARCHAR(191) DEFAULT \"\", `created_at` TIMESTAMP NULL DEFAULT NOW(), `updated_at` TIMESTAMP NULL DEFAULT NOW(), PRIMARY KEY (`id`), UNIQUE KEY `projects_alt_unique_idx` (`title`, `organization_name`, `description`, `pass_type`) );"
	tempUserProjectsTable  = "CREATE TEMPORARY TABLE `user_projects` ( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `user_id` int(10) unsigned NOT NULL, `project_id` int(10) unsigned NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `projects_user_and_project_unique_idx` (`user_id`, `project_id`) );"
)

const (
	insertPassesTable        = "INSERT INTO `passes` (`serial_number`, `authentication_token`, `pass_type_id`, `updated_at`) VALUES ('%s', '%s', '%s', '%s');"
	insertRegistrationsTable = "INSERT INTO `registrations` (`device_id`, `push_token`, `serial_number`, `pass_type_id`) VALUES ('%s', '%s', '%s', '%s');"
	insertUsersTable         = "INSERT INTO `users` (`email`, `username`, `password_hash`) VALUES ('%s', '%s', '%s');"
	insertProjectsTable      = "INSERT INTO `projects` (`organization_id`, `description`, `pass_type`) VALUES ('%s', '%s', '%s');"
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

func fakeUsername() string {
	return uuid.NewV4().String()
}

func fakeEmail() string {
	return uuid.NewV4().String() + "@example.com"
}

func fakeString() string {
	return uuid.NewV4().String()
}
