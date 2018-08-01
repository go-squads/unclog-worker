DROP TABLE IF EXISTS log_metrics;

CREATE TABLE log_metrics (
    id serial primary key,
    "timestamp" timestamp without time zone,
    log_level text not null,
    log_level_id integer not null,
    quantity integer not null
);
