-- file_msg_data 表
CREATE TABLE IF NOT EXISTS file_msg_data (
    msg_id BIGINT UNSIGNED PRIMARY KEY,
    file_hash VARCHAR(64) NOT NULL,
    file_intlID INT UNSIGNED UNIQUE,
    file_size INT UNSIGNED,
    file_orglName VARCHAR(192),
    file_msgTime TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6), -- 写入时的时间
    file_realName VARCHAR(320),
    file_first32 INT UNSIGNED
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- msg 表
CREATE TABLE IF NOT EXISTS msg (
    msg_id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    group_id INT UNSIGNED,
    msg_content TEXT,
    msg_msgTime TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
    msg_uid INT UNSIGNED,
    msg_fileHash VARCHAR(128),
    msg_type SMALLINT UNSIGNED,
    INDEX idx_msg_time (msg_msgTime DESC), -- 时间索引
    INDEX idx_group (group_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci
PARTITION BY
    RANGE COLUMNS (msg_id) (
        PARTITION p0
        VALUES
            LESS THAN (12000000), -- 每个分区最大1200万条
            PARTITION p1
        VALUES
            LESS THAN (24000000),
            PARTITION perror
        VALUES
            LESS THAN (MAXVALUE)
    );