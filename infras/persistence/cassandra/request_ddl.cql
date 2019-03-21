-- CREATE KEYSPACE crawler WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1 };
-- use crawler;

CREATE TABLE request
(
  namespace    text,
  url          text,
  next_request bigint,
  method       text,
  body         blob,
  cookie       blob,
  job_status   int,
  last_request bigint,
  stats        map<text, double>,
  PRIMARY KEY ((namespace, url), next_request)
) WITH CLUSTERING ORDER BY (next_request ASC);
