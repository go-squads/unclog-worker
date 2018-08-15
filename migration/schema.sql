DROP TABLE IF EXISTS log_metrics_v1, log_metrics_v2;

CREATE TABLE log_metrics_v1 (
    id serial primary key,
    "timestamp" timestamp without time zone,
    log_level text not null,
    quantity integer not null
);

CREATE TABLE log_metrics_v2 (
    id serial primary key,
    "timestamp" timestamp without time zone,
    app_name text not null,
    node_id text not null,
    log_level text not null,
    quantity integer not null
);

create table alerts (
  id serial primary key,
  app_name text not null,
  log_level text not null,
  duration integer not null,
  "limit" integer not null,
  callback text not null
)