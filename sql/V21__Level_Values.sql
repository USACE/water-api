-- level_value definition
CREATE TABLE IF NOT EXISTS level_value (
    level_id UUID NOT NULL REFERENCES level(id),
    day_of_year INTEGER NOT NULL CHECK (day_of_year > 0 AND day_of_year < 366),
    value DOUBLE PRECISION NOT NULL,
    CONSTRAINT unique_level_day UNIQUE(level_id, day_of_year)
);

-- INSERT INTO level_value (level_id, day_of_year, value) VALUES
--     ('a838d96f-2e11-40f8-a699-4e81f9b79a09', 1, 664),
--     ('a838d96f-2e11-40f8-a699-4e81f9b79a09', 366, 664);

-- grant read
GRANT SELECT ON
    level_value
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    level_value
TO water_writer;
