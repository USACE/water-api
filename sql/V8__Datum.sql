-- Create vertical_datum table
CREATE TABLE IF NOT EXISTS vertical_datum (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR UNIQUE NOT NULL
);

-- Grant read
GRANT SELECT ON
    vertical_datum
TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON
    vertical_datum
TO water_writer;

INSERT into vertical_datum (id, name) VALUES
(0, 'UNKNOWN'),
(1, 'COE1912'),
(2, 'NGVD29'),
(3, 'NAVD88');