-- CREATE KEYSPACE crawler WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1 };
-- use crawler;

CREATE TABLE request_pq
(
  namespace    text,
  next_request bigint,
  url          text,
  primary key ( (namespace), next_request, url)
) WITH CLUSTERING ORDER BY (next_request ASC);

CREATE TABLE request
(
  namespace    text,
  url          text,
  method       text,
  body         blob,
  cookie       blob,
  job_status   int,
  last_request bigint,
  next_request bigint,
  stats        blob,
  PRIMARY KEY ((namespace), url)
) WITH CLUSTERING ORDER BY (url ASC);
