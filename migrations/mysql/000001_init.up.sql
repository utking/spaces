CREATE TABLE IF NOT EXISTS `user` (
    id varchar(36) DEFAULT (UUID()) PRIMARY KEY,
    username VARCHAR(16) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    auth_key VARCHAR(32),
    account_activation_token VARCHAR(255) DEFAULT NULL,
    status SMALLINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (status),
    UNIQUE (username),
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS `auth_item` (
  `name` varchar(64) NOT NULL,
  `type` smallint NOT NULL,
  `description` text,
  `rule_name` varchar(64) DEFAULT NULL,
  `data` blob,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`name`),
  KEY `rule_name` (`rule_name`),
  KEY `idx-auth_item-type` (`type`)
);

CREATE TABLE IF NOT EXISTS `auth_assignment` (
  `item_name` varchar(64) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`item_name`,`user_id`),
  KEY `idx-auth_assignment-user_id` (`user_id`),
  CONSTRAINT `auth_assignment_ibfk_1` FOREIGN KEY (`item_name`) REFERENCES `auth_item` (`name`) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO `user` (`id`, `username`, `password_hash`, `email`, `auth_key`, `status`)
  VALUES ('efb9e25d-4323-11f0-a5be-0242ac160002', 'root', '$2a$13$M6tMq/VVqDKsYaHYm149Q.IreL1q4KCLZy1SIo8zxmhGdVxrE4sa.',
         'user@gmail.com', 'QvV3Y5a8fClIqOwT1Y7eBjFnKtQyW5b-', 10);

INSERT INTO `auth_item` (`name`, `type`, `description`, `rule_name`, `data`)
  VALUES ('admin', 1, 'Admin role', NULL, NULL),
         ('user', 1, 'User role', NULL, NULL);

INSERT INTO `auth_assignment` (`item_name`, `user_id`)
  VALUES ('admin', 'efb9e25d-4323-11f0-a5be-0242ac160002');

CREATE TABLE IF NOT EXISTS `note` (
    id varchar(36) DEFAULT (UUID()) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    title VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    `tags` JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (`title`, `user_id`),
    INDEX note_tags_idx ( (cast(`tags` as char(32) array)) ),
    INDEX user_id_idx (user_id),
    FOREIGN KEY (user_id) REFERENCES `user` (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS `password_record` (
    id varchar(36) DEFAULT (UUID()) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `url` VARCHAR(4096) DEFAULT '',
    `description` TEXT,
    `tags` JSON DEFAULT NULL,
    `username` VARBINARY(1024) NOT NULL,
    `secret` VARBINARY(4096) NOT NULL,
    INDEX passwd_tags_idx ( (cast(`tags` as char(32) array)) ),
    INDEX user_id_idx (user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (`name`, `user_id`),
    FOREIGN KEY (user_id) REFERENCES `user` (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS `bookmark` (
    id varchar(36) DEFAULT (UUID()) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `url` VARCHAR(4096) NOT NULL,
    `tags` JSON DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX tags_idx ( (cast(`tags` as char(32) array)) ),
    INDEX user_id_idx (user_id),
    FOREIGN KEY (user_id) REFERENCES `user` (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS `last_opened` (
    `user_id` varchar(36) NOT NULL,
    `item_type` ENUM('note_id', 'bookmark_tag') NOT NULL,
    `item_id` varchar(36) NOT NULL,
    INDEX item_type_idx (item_type),
    INDEX item_id_idx (item_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, item_type)
);
