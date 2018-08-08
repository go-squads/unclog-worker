DROP TABLE IF EXISTS log_metrics;

CREATE TABLE log_metrics (
    id serial primary key,
    "timestamp" timestamp without time zone,
    app_name text not null,
    node_id text not null,
    log_level text not null,
    quantity integer not null
);
