-- provider table
-- setup manualy to enforce contraints and defaults
CREATE TABLE IF NOT EXISTS provider (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL,
    slug VARCHAR UNIQUE NOT NULL,
    parent_id UUID references provider(id)
);

-- copy office table into 
INSERT into provider (select id, name, lower(symbol), parent_id from office);
