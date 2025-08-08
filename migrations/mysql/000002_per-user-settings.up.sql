CREATE TABLE IF NOT EXISTS `user_settings` (
  user_id varchar(36) PRIMARY KEY,
  `value` JSON NOT NULL DEFAULT ( '{}' ),
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
);