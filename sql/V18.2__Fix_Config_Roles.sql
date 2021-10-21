-- Grant read
GRANT SELECT ON
    config
TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON
    config
TO water_writer;