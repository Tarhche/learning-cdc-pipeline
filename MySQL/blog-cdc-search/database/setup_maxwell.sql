-- Create Maxwell user for CDC
CREATE USER IF NOT EXISTS 'maxwell'@'%' IDENTIFIED BY 'maxwell123';

-- Grant Maxwell database permissions (for storing state)
GRANT ALL ON maxwell.* TO 'maxwell'@'%';

-- Grant replication permissions to all database tables for CDC
GRANT SELECT, REPLICATION CLIENT, REPLICATION SLAVE ON *.* TO 'maxwell'@'%';

-- Maxwell schema (required for Maxwell to track its state)
CREATE DATABASE IF NOT EXISTS maxwell;

-- Flush privileges
FLUSH PRIVILEGES;
