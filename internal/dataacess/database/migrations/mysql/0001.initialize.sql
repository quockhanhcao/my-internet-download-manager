CREATE TABLE IF NOT EXISTS `accounts` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `account_name` VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS `account_passwords` (
  `of_account_id` BIGINT UNSIGNED PRIMARY KEY,
  `hash` VARCHAR(128) NOT NULL,
  FOREIGN KEY (`of_account_id`) REFERENCES `accounts`(`account_id`)
);

CREATE TABLE IF NOT EXISTS `download_tasks` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `of_account_id` BIGINT UNSIGNED,
  `download_type` SMALLINT NOT NULL,
  `url` TEXT NOT NULL,
  `download_status` SMALLINT NOT NULL,
  `metadata` TEXT NOT NULL,
  FOREIGN KEY (`of_account_id`) REFERENCES `accounts`(`account_id`)
);

CREATE TABLE IF NOT EXISTS `token_public_keys` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `public_key` VARBINARY(4096) NOT NULL
)
