-- CREATE KEYSPACE crawler WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1 };
-- use crawler;

CREATE TABLE request
(
  url_hash     text,
  next_request bigint,
  url          text, -- url string
  method       text,
  body         blob,
  cookie       blob,
  job_status   int,
  last_request bigint,
  stats        map<text, double>,
  PRIMARY KEY ((url_hash), next_request)
) WITH CLUSTERING ORDER BY (next_request DESC);