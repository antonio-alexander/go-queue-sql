-- DROP DATABASE go_queue_sql;
CREATE DATABASE IF NOT EXISTS go_queue_sql;

-- DROP TABLE queue;
CREATE TABLE IF NOT EXISTS queue (
	id SERIAL,
	data BLOB
) ENGINE = InnoDB;