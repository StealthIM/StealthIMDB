CREATE TABLE IF NOT EXISTS `groups` (
    `groupid` INT(32) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `password` VARCHAR(128),
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `owner_uid` INT(32) NOT NULL,
    `config_bitset` INT(64) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `group_user_table` (
    `groupid` INT(32) NOT NULL,
    `username` VARCHAR(64) NOT NULL,
    `type` ENUM('member', 'manager', 'owner', 'other') NOT NULL,
    `config_bitset` INT(64) NOT NULL DEFAULT 0,
    PRIMARY KEY (`groupid`, `username`)
);

CREATE TABLE IF NOT EXISTS `group_message_table` (
    `groupid` INT(32) NOT NULL,
    `username` VARCHAR(64) NOT NULL,
    `msgid` BIGINT(64) NOT NULL,
    `send_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`groupid`, `msgid`)
);