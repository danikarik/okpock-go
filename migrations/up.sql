CREATE TABLE IF NOT EXISTS `passes` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `authentication_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `passes_serial_number_unique` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `registrations` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `uuid` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `device_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `push_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `registrations_serial_number_unique` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `logs` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `uuid` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `remote_address` VARCHAR(191) COLLATE utf8mb4_unicode_ci NULL,
    `request_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `logs_uuid_unique` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `users` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `role` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "client",
    `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password_hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `confirmed_at` TIMESTAMP NULL DEFAULT NULL,
    `invited_at` TIMESTAMP NULL DEFAULT NULL,
    `confirmation_token` VARCHAR(191) DEFAULT NULL,
    `confirmation_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `recovery_token` VARCHAR(191) DEFAULT NULL,
    `recovery_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `email_change_token` VARCHAR(191) DEFAULT NULL,
    `email_change` VARCHAR(191) DEFAULT NULL,
    `email_change_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `last_signin_at` TIMESTAMP NULL DEFAULT NULL,
    `raw_app_metadata` text DEFAULT NULL,
    `raw_user_metadata` text DEFAULT NULL,
    `is_super_admin` TINYINT(1) DEFAULT 0,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_username_unique` (`username`),
    UNIQUE KEY `users_email_unique` (`email`),
    UNIQUE KEY `users_username_and_email_unique` (`username`, `email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `organizations` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` INT(10) unsigned NOT NULL,
    `title` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `raw_metadata` text DEFAULT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    FOREIGN KEY `organizations_user_reference` (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    UNIQUE KEY `organizations_title_and_user_unique` (`user_id`, `title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `projects` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `organization_id` INT(10) unsigned NOT NULL,
    `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `background_image` VARCHAR(191) DEFAULT NULL,
    `footer_image` VARCHAR(191) DEFAULT NULL,
    `icon_image` VARCHAR(191) DEFAULT NULL,
    `strip_image` VARCHAR(191) DEFAULT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    FOREIGN KEY `projects_organization_reference` (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE,
    UNIQUE KEY `projects_title_and_user_unique` (`organization_id`, `description`, `pass_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
