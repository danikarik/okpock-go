CREATE TABLE IF NOT EXISTS `passes` (
    `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL,
    `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `authentication_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `passes_serial_number_unique_idx` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `registrations` (
    `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL,
    `device_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `push_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `registrations_serial_number_unique_idx` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `logs` (
    `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL,
    `remote_address` VARCHAR(191) COLLATE utf8mb4_unicode_ci NULL,
    `request_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL,
    `role` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "client",
    `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password_hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `confirmed_at` TIMESTAMP NULL DEFAULT NULL,
    `invited_at` TIMESTAMP NULL DEFAULT NULL,
    `confirmation_token` VARCHAR(191) DEFAULT "",
    `confirmation_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `recovery_token` VARCHAR(191) DEFAULT "",
    `recovery_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `email_change_token` VARCHAR(191) DEFAULT "",
    `email_change` VARCHAR(191) DEFAULT "",
    `email_change_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `last_signin_at` TIMESTAMP NULL DEFAULT NULL,
    `raw_app_metadata` text DEFAULT NULL,
    `raw_user_metadata` text DEFAULT NULL,
    `is_super_admin` TINYINT(1) DEFAULT 0,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_username_unique_idx` (`username`),
    UNIQUE KEY `users_email_unique_idx` (`email`),
    UNIQUE KEY `users_username_and_email_unique_idx` (`username`, `email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `projects` (
    `id` VARCHAR(144) COLLATE utf8mb4_unicode_ci NOT NULL,
    `title` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `organization_name` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `background_image` VARCHAR(191) DEFAULT "",
    `footer_image` VARCHAR(191) DEFAULT "",
    `icon_image` VARCHAR(191) DEFAULT "",
    `strip_image` VARCHAR(191) DEFAULT "",
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `projects_alt_unique_idx` (`title`, `organization_name`, `description`, `pass_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
