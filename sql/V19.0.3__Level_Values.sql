-- level_value definition
CREATE TABLE IF NOT EXISTS level_value (
    level_id UUID NOT NULL REFERENCES level(id),
    julian_date DOUBLE PRECISION NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    CONSTRAINT unique_level_julian UNIQUE(level_id, julian_date)
);

-- grant read
GRANT SELECT ON
    level_value
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    level_value
TO water_writer;
