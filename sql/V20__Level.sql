-- level definition
CREATE TABLE IF NOT EXISTS level (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES location(id),
    level_kind_id UUID NOT NULL REFERENCES level_kind(id),
    CONSTRAINT unique_location_level_kind UNIQUE(location_id, level_kind_id)
);

-- INSERT INTO level (id, location_id, level_kind_id) VALUES
--     ('a838d96f-2e11-40f8-a699-4e81f9b79a09', '78d77de9-e4cc-486b-ab41-5d0d93771c4c', '43e6ecff-32d0-4e03-ba79-f05a9ed5924d');

-- grant read
GRANT SELECT ON
    level
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    level
TO water_writer;
