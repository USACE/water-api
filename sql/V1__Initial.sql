CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS location_kind (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS location (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    office_id UUID,
    name VARCHAR,
    public_name VARCHAR,
    slug VARCHAR UNIQUE NOT NULL,
    geometry geometry,
    kind_id UUID NOT NULL REFERENCES location_kind (id),
    creator UUID,
    create_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    updater UUID,
    update_date TIMESTAMPTZ,
    CONSTRAINT office_unique_name UNIQUE(office_id,name)
);
