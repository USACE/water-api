-- level_kind definition
CREATE TABLE IF NOT EXISTS level_kind (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR UNIQUE NOT NULL
);

-- grant read
GRANT SELECT ON
    level_kind
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    level_kind
TO water_writer;
