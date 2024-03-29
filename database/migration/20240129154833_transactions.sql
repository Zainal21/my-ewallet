-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `order_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ref_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gross_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `piece` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `amount` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `note` Text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` enum('process','success','failed') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'process',
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ledgers_unique` (`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
