-- MYSQL 
-- MySQLでvarcharのインデックスできるバイト数を3072にするために以下の変更を行う
-- MySQL 5.5以下を利用している場合は SET GLOBAL innodb_file_per_table=1;
-- MySQL 5.6以下を利用している場合は SET GLOBAL innodb_large_prefix=1;

CREATE TABLE response_meta
(
  namespace varchar(64),
  `date`     int,
  response_hash varchar(256),
  url        varchar(1024),
  primary key (namespace, `date`)
  )
ROW_FORMAT=COMPRESSED
ENGINE=INNODB
DEFAULT CHARSET=utf8mb4
PARTITION BY key(namespace, `date`)
PARTITIONS 100
;
