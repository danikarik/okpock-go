CREATE TABLE IF NOT EXISTS `passes` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `authentication_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `passes_serial_number_unique_idx` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `registrations` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `device_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `push_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `serial_number` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `registrations_serial_number_unique_idx` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `logs` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `remote_address` VARCHAR(191) COLLATE utf8mb4_unicode_ci NULL,
    `request_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `message` TEXT NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `users` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `role` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "client",
    `username` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password_hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `confirmed_at` TIMESTAMP NULL DEFAULT NULL,
    `invited_at` TIMESTAMP NULL DEFAULT NULL,
    `confirmation_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `confirmation_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `recovery_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `recovery_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `email_change_token` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `email_change` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `email_change_sent_at` TIMESTAMP NULL DEFAULT NULL,
    `last_signin_at` TIMESTAMP NULL DEFAULT NULL,
    `raw_app_metadata` TEXT DEFAULT NULL,
    `raw_user_metadata` TEXT DEFAULT NULL,
    `is_super_admin` TINYINT(1) DEFAULT 0,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_username_unique_idx` (`username`),
    UNIQUE KEY `users_email_unique_idx` (`email`),
    UNIQUE KEY `users_username_and_email_unique_idx` (`username`, `email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `projects` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `organization_name` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `pass_type` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `background_image` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `background_image_2x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `background_image_3x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `footer_image` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `footer_image_2x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `footer_image_3x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `icon_image` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `icon_image_2x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `icon_image_3x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `logo_image` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `logo_image_2x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `logo_image_3x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `strip_image` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `strip_image_2x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `strip_image_3x` VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT "",
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `projects_alt_unique_idx` (`title`, `organization_name`, `description`, `pass_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_projects` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` INT(10) unsigned NOT NULL,
    `project_id` INT(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `projects_user_and_project_unique_idx` (`user_id`, `project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `uploads` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `uuid` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `filename` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `hash` VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uploads_alt_unique_idx` (`uuid`, `filename`, `hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_uploads` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` INT(10) unsigned NOT NULL,
    `upload_id` INT(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uploads_user_and_upload_unique_idx` (`user_id`, `upload_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_pass_cards` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `project_id` INT(10) unsigned NOT NULL,
    `pass_card_id` INT(10) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pass_cards` (
    `id` INT(10) unsigned NOT NULL AUTO_INCREMENT,
    `raw_data` TEXT DEFAULT NULL,
    `created_at` TIMESTAMP NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL DEFAULT NOW(),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
