-- level definition
CREATE TABLE IF NOT EXISTS level (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    location_id UUID,
    level_kind_id UUID NOT NULL REFERENCES level_kind(id),
    CONSTRAINT unique_location_level_kind UNIQUE(location_id, level_kind_id)
);

-- grant read
GRANT SELECT ON
    level
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    level
TO water_writer;

