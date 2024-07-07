-- DROP DATABASE go_queue_sql;
CREATE DATABASE IF NOT EXISTS go_queue_sql;

-- DROP TABLE queue;
CREATE TABLE IF NOT EXISTS queue (
	id SERIAL,
	priority INT DEFAULT 0,
	data BLOB
) ENGINE = InnoDB;