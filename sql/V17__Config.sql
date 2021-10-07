-- config (application config variables)
CREATE TABLE IF NOT EXISTS config (
    config_name VARCHAR UNIQUE NOT NULL,
    config_value VARCHAR NOT NULL
);

INSERT INTO config (config_name, config_value) VALUES
('write_to_bucket', 'castle-data-develop');