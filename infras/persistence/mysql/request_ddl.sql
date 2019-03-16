CREATE TABLE request
(
  namespace    varchar(64),
  url          text, -- url string
  url_hash     char(32),
  method       varchar(10),
  body         blob,
  cookie       blob,
  job_status   int,
  next_request bigint,
  last_request bigint,
  stats        blob,
  primary key (namespace, url_hash),
  index (next_request)
)
  ROW_FORMAT = COMPRESSED
  ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  PARTITION BY key (namespace) PARTITIONS 100
;
