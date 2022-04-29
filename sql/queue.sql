
-- DROP DATABASE IF EXISTS go-queue-sql;
CREATE DATABASE IF NOT EXISTS go_queue_sql;

USE go_queue_sql;

-- DROP TABLE IF EXISTS queue
CREATE TABLE IF NOT EXISTS queue (
    id SERIAL,
    data BLOB
) ENGINE = InnoDB;