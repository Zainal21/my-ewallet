-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `ledgers` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ref_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `current_deposit` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `change_deposit` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `final_deposit` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ledgers_unique` (`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ledgers;
-- +goose StatementEnd
