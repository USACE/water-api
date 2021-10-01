-- level_kind definition
CREATE TABLE IF NOT EXISTS level_kind (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR UNIQUE NOT NULL
);

INSERT INTO level_kind (id, slug, name) VALUES
    ('43e6ecff-32d0-4e03-ba79-f05a9ed5924d', 'top-of-flood-control', 'Top of Flood Control'),
    ('7a998105-2d91-4b2d-ab5d-9d2fe12b9125', 'bottom-of-flood-control', 'Bottom of Flood Control');


-- grant read
GRANT SELECT ON
    level_kind
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    level_kind
TO water_writer;
