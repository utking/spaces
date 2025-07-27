CREATE TABLE `user` (
  `id` varchar(36) PRIMARY KEY,
  `username` varchar(16) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `auth_key` varchar(32),
  `account_activation_token` varchar(255) DEFAULT NULL,
  `status` SMALLINT DEFAULT 0 NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT current_timestamp,
  `updated_at` DATETIME NOT NULL DEFAULT current_timestamp
);

CREATE UNIQUE INDEX `UQE_user_username` ON `user` (`username`);
CREATE UNIQUE INDEX `UQE_user_email` ON `user` (`email`);
CREATE INDEX `idx_user_status` ON `user` (`status`);

CREATE TABLE IF NOT EXISTS `auth_item` (
  `name` varchar(64) PRIMARY KEY,
  `type` smallint NOT NULL, 
  `description` text,
  `rule_name` varchar(64) DEFAULT NULL,
  `data` blob,
  `created_at` DATETIME NOT NULL DEFAULT current_timestamp,
  `updated_at` DATETIME NOT NULL DEFAULT current_timestamp
);

CREATE UNIQUE INDEX rule_name ON `auth_item` (`rule_name`);
CREATE INDEX idx_auth_item_type ON `auth_item` (`type`);

CREATE TABLE IF NOT EXISTS `auth_assignment` (
  `item_name` varchar(64) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT current_timestamp,
  PRIMARY KEY (`item_name`,`user_id`),
  FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  FOREIGN KEY (`item_name`) REFERENCES `auth_item` (`name`)
);

CREATE INDEX idx_auth_assignment_user_id ON `auth_assignment` (`user_id`);

INSERT INTO `user` (`id`, `username`, `password_hash`, `email`, `auth_key`, `status`)
  VALUES ('efb9e25d-4323-11f0-a5be-0242ac160002', 'root', '$2a$13$M6tMq/VVqDKsYaHYm149Q.IreL1q4KCLZy1SIo8zxmhGdVxrE4sa.',
         'user@gmail.com', 'QvV3Y5a8fClIqOwT1Y7eBjFnKtQyW5b-', 10);

INSERT INTO `auth_item` (`name`, `type`, `description`, `rule_name`, `data`)
  VALUES ('admin', 1, 'Admin role', NULL, NULL),
         ('user', 1, 'User role', NULL, NULL);

INSERT INTO `auth_assignment` (`item_name`, `user_id`)
  VALUES ('admin', 'efb9e25d-4323-11f0-a5be-0242ac160002');

CREATE TABLE IF NOT EXISTS `note` (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    title VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    `tags` TEXT DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT current_timestamp,
    updated_at DATETIME NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES `user` (id)
);

CREATE INDEX idx_note_user_id ON `note` (user_id);
CREATE UNIQUE INDEX idx_note_title_user ON `note` (`title`, `user_id`);

CREATE TABLE IF NOT EXISTS `password_record` (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `username` blob NOT NULL,
    `url` VARCHAR(4096) DEFAULT '',
    `description` TEXT,
    `tags` TEXT DEFAULT NULL,
    `secret` blob NOT NULL,
    created_at DATETIME NOT NULL DEFAULT current_timestamp,
    updated_at DATETIME NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES `user` (id)
);

CREATE INDEX idx_password_record_user_id ON `password_record` (user_id);
CREATE UNIQUE INDEX idx_password_record_name_user ON `password_record` (`name`, `user_id`);

CREATE TABLE IF NOT EXISTS `bookmark` (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `url` VARCHAR(4096) NOT NULL,
    `tags` TEXT DEFAULT NULL,
    `created_at` DATETIME NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES `user` (id)
);

CREATE INDEX idx_bookmark_user_id ON `bookmark` (user_id);

CREATE TABLE IF NOT EXISTS `last_opened` (
    `user_id` varchar(36) NOT NULL,
    `item_type` CHECK (item_type IN ('note_id', 'bookmark_tag')) NOT NULL,
    `item_id` varchar(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id),
    PRIMARY KEY (user_id, item_type)
);

CREATE INDEX idx_last_opened_user_id ON `last_opened` (user_id);
CREATE INDEX idx_last_opened_item_type ON `last_opened` (item_type);
CREATE INDEX idx_last_opened_item_id ON `last_opened` (item_id);