CREATE TABLE response
(
  response_hash varchar(256) primary key, -- sha512 hashed body
  url           varchar(1024),            -- url string
  body          LONGTEXT,
  `date`        int,
  created_at    datetime,
  index(`date`)
)
  ROW_FORMAT = COMPRESSED
  ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  PARTITION BY key (response_hash)
    PARTITIONS 100
;
