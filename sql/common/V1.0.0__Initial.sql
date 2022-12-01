CREATE extension IF NOT EXISTS "uuid-ossp";

-- config (application config variables)
CREATE TABLE IF NOT EXISTS config (
    config_name VARCHAR UNIQUE NOT NULL,
    config_value VARCHAR NOT NULL
);

------------
-- PROVIDER
------------

CREATE TABLE IF NOT EXISTS provider (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL,
    slug VARCHAR UNIQUE NOT NULL,
    parent_id UUID references provider(id)
);

----------------
-- DATASOURCE
----------------

-- datatype table
CREATE TABLE IF NOT EXISTS datatype (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    uri VARCHAR
);

-- datasource table
CREATE TABLE IF NOT EXISTS datasource (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES provider(id),
    datatype_id UUID NOT NULL REFERENCES datatype(id),
    CONSTRAINT datasource_unique_provider_datatype UNIQUE(provider_id, datatype_id)
);

------------
-- LOCATION
------------

CREATE TABLE IF NOT EXISTS location (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    datasource_id UUID NOT NULL REFERENCES datasource(id),
    slug VARCHAR UNIQUE NOT NULL,
    code VARCHAR NOT NULL,
    geometry geometry NOT NULL, 
    state_id INTEGER REFERENCES tiger_data.state_all(gid),
    create_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_date TIMESTAMPTZ,
    attributes JSONB NOT NULL DEFAULT '{}'::jsonb
);
-- Ensure case-insensitive uniqueness in code column (within datasource context)
CREATE UNIQUE INDEX location_case_insensitive_unique_code ON location (datasource_id, LOWER(code));


--------------
-- TIMESERIES
--------------

-- timeseries table
CREATE TABLE IF NOT EXISTS timeseries (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    datasource_id UUID NOT NULL REFERENCES datasource(id),
    datasource_key VARCHAR NOT NULL,
    location_id UUID NOT NULL REFERENCES location(id),
    latest_value DOUBLE PRECISION,
    latest_time TIMESTAMPTZ,
    etl_values_enabled BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT timeseries_unique_datasource UNIQUE(datasource_id, datasource_key)
);
-- Ensure case-insensitive uniqueness in datasource_key column (within a single datasource)
CREATE UNIQUE INDEX timeseries_case_insensitive_unique_datasource_key ON timeseries (datasource_id, LOWER(datasource_key));

-- timeseries_value
CREATE TABLE IF NOT EXISTS timeseries_value (
    timeseries_id UUID NOT NULL REFERENCES timeseries (id) ON DELETE CASCADE,
    time TIMESTAMPTZ NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    CONSTRAINT timeseries_id_unique_time UNIQUE(timeseries_id, time),
    PRIMARY KEY (timeseries_id, time)
);

-------------------
-- TIMESERIES_GROUP
-------------------
CREATE TABLE IF NOT EXISTS timeseries_group (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    provider_id UUID NOT NULL REFERENCES provider(id),
    CONSTRAINT provider_unique_timeseries_group_name UNIQUE(provider_id, name)
);

CREATE TABLE IF NOT EXISTS timeseries_group_members (
    timeseries_group_id UUID NOT NULL REFERENCES timeseries_group(id) ON DELETE CASCADE,
    timeseries_id UUID NOT NULL REFERENCES timeseries(id) ON DELETE CASCADE,
    CONSTRAINT timeseries_group_unique_timeseries UNIQUE(timeseries_group_id, timeseries_id)
);


--------------
-- CHART
--------------

-- chart table
CREATE TABLE IF NOT EXISTS chart (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    location_id UUID REFERENCES location(id),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    provider_id UUID NOT NULL REFERENCES provider(id),
    CONSTRAINT provider_unique_chart_name UNIQUE(provider_id, name)
);

-- chart_variable_mapping
CREATE TABLE IF NOT EXISTS chart_variable_mapping (
    chart_id UUID NOT NULL REFERENCES chart(id) ON DELETE CASCADE,
    variable VARCHAR NOT NULL,
    timeseries_id UUID NOT NULL REFERENCES timeseries(id) ON DELETE CASCADE,
    CONSTRAINT chart_unique_variable UNIQUE(chart_id, variable)
);
