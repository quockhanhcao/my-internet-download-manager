CREATE TABLE IF NOT EXISTS `users` (
  `user_id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS `user_passwords` (
  `of_user_id` BIGINT UNSIGNED PRIMARY KEY,
  `hash` VARCHAR(128) NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`)
);

CREATE TABLE IF NOT EXISTS `download_tasks` (
  `of_user_id` BIGINT UNSIGNED PRIMARY KEY,
  `download_type` SMALLINT NOT NULL,
  `url` TEXT NOT NULL,
  `download_status` SMALLINT NOT NULL,
  `metadata` TEXT NOT NULL,
  FOREIGN KEY (`of_user_id`) REFERENCES `users`(`user_id`)
);
