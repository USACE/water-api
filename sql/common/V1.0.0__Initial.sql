CREATE extension IF NOT EXISTS "uuid-ossp";

-- config (application config variables)
CREATE TABLE IF NOT EXISTS config (
    config_name VARCHAR UNIQUE NOT NULL,
    config_value VARCHAR NOT NULL
);

------------------------
-- CWMS LOCATION KIND
------------------------

CREATE TABLE IF NOT EXISTS cwms_location_kind (
    name VARCHAR UNIQUE NOT NULL
);

-------------------
-- USGS SITE TYPE
-------------------

CREATE TABLE IF NOT EXISTS usgs_site_type (
    abbreviation VARCHAR UNIQUE NOT NULL,
    name VARCHAR UNIQUE NOT NULL
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
    -- CONSTRAINT datatype_unique_uri UNIQUE(slug, uri)  // todo; delete; this constraint will always be met because slug is globally unique
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



-- Create usgs_parameter table
CREATE TABLE IF NOT EXISTS usgs_parameter (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,
    description VARCHAR NOT NULL
);


-- Create usgs_site_parameters table
CREATE TABLE IF NOT EXISTS usgs_site_parameters (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES location(id),
    parameter_id UUID NOT NULL REFERENCES usgs_parameter(id),
    CONSTRAINT site_unique_param UNIQUE(location_id, parameter_id)
);

-- usgs_measurements
CREATE TABLE IF NOT EXISTS usgs_measurements (
    time TIMESTAMPTZ NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    usgs_site_parameters_id UUID NOT NULL REFERENCES usgs_site_parameters (id) ON DELETE CASCADE,
    CONSTRAINT site_parameters_unique_time UNIQUE(usgs_site_parameters_id, time),
    PRIMARY KEY (usgs_site_parameters_id, time)
);

-- -- watershed
-- CREATE TABLE IF NOT EXISTS watershed (
--     id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
--     slug VARCHAR UNIQUE NOT NULL,
--     name VARCHAR NOT NULL,
--     geometry geometry NOT NULL DEFAULT ST_GeomFromText('POLYGON ((0 0, 0 0, 0 0, 0 0, 0 0))',4326),
--     provider_id UUID NOT NULL REFERENCES provider(id),
-- 	deleted boolean NOT NULL DEFAULT false
-- );

-- -- watershed_usgs_sites
-- CREATE TABLE IF NOT EXISTS watershed_usgs_sites (
--     watershed_id UUID NOT NULL REFERENCES watershed(id),
--     usgs_site_parameter_id UUID NOT NULL REFERENCES usgs_site_parameters(id),
--     CONSTRAINT watershed_unique_site_param UNIQUE(watershed_id, usgs_site_parameter_id)
-- );
-- -- Add comment to describe table
-- COMMENT ON TABLE watershed_usgs_sites IS 'This is a bridge table.  Each entry represent a watershed/site/parameter that has been requested for USGS data acquisition.';


-------------------
-- SHAPEFILE UPLOAD
-------------------

-- upload_status definition
CREATE TABLE IF NOT EXISTS upload_status (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL
);

-- -- watershed_shapefile_uploads definition
-- CREATE TABLE IF NOT EXISTS watershed_shapefile_uploads (
-- 	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
-- 	watershed_id UUID NOT NULL REFERENCES watershed(id),
-- 	file VARCHAR NOT NULL,
-- 	date_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
-- 	file_size INTEGER NOT NULL,
--     processing_info VARCHAR,
--     user_id UUID,
--     upload_status_id UUID NOT NULL DEFAULT 'b5d777fc-c46b-4a10-a488-1415e3d7849d' REFERENCES upload_status(id)
-- );

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

-- timeseries_measurement
CREATE TABLE IF NOT EXISTS timeseries_measurement (
    timeseries_id UUID NOT NULL REFERENCES timeseries (id) ON DELETE CASCADE,
    time TIMESTAMPTZ NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    CONSTRAINT timeseries_id_unique_time UNIQUE(timeseries_id, time),
    PRIMARY KEY (timeseries_id, time)
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
    type_id UUID NOT NULL,
    provider_id UUID NOT NULL REFERENCES provider(id),
    CONSTRAINT provider_unique_name UNIQUE(provider_id, name)
);

-- chart_variable_mapping
CREATE TABLE IF NOT EXISTS chart_variable_mapping (
    chart_id UUID NOT NULL REFERENCES chart(id),
    variable VARCHAR NOT NULL,
    timeseries_id UUID NOT NULL REFERENCES timeseries(id) ON DELETE CASCADE,
    CONSTRAINT chart_id_unique_variable UNIQUE(chart_id, variable)
);