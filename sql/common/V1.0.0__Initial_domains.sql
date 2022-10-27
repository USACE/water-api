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
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR UNIQUE NOT NULL
);

-------------------
-- USGS SITE TYPE
-------------------

CREATE TABLE IF NOT EXISTS usgs_site_type (
    id UUID PRIMARY KEY NOT NULL,
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

-- datasource_type table
CREATE TABLE IF NOT EXISTS datasource_type (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    uri VARCHAR NOT NULL,
    CONSTRAINT datasource_type_slug_unique_uri UNIQUE(slug, uri)
);

-- datasource table
CREATE TABLE IF NOT EXISTS datasource (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES provider(id),
    datasource_type_id UUID NOT NULL REFERENCES datasource_type(id),
    CONSTRAINT datasource_unique_provider_id UNIQUE(provider_id, datasource_type_id)
);

------------
-- LOCATION
------------

CREATE TABLE IF NOT EXISTS location (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    datasource_id UUID NOT NULL REFERENCES datasource(id),
    slug VARCHAR UNIQUE NOT NULL,
    geometry geometry NOT NULL, 
    state_id INTEGER REFERENCES tiger_data.state_all(gid),
    create_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_date TIMESTAMPTZ
);

----------------
-- CWMS LOCATION
----------------

CREATE TABLE IF NOT EXISTS cwms_location (
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    name VARCHAR,
    public_name VARCHAR,
    kind_id UUID NOT NULL REFERENCES cwms_location_kind(id)
);

----------------
-- USGS SITE
----------------

CREATE TABLE IF NOT EXISTS usgs_site (
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    site_number VARCHAR UNIQUE NOT NULL,
    station_name VARCHAR NOT NULL,
    site_type_id UUID NOT NULL REFERENCES usgs_site_type(id)
);


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

----------------
-- NWS SITE
----------------

CREATE TABLE IF NOT EXISTS nws_site (
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    nws_li VARCHAR UNIQUE NOT NULL
);

-- watershed
CREATE TABLE IF NOT EXISTS watershed (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    geometry geometry NOT NULL DEFAULT ST_GeomFromText('POLYGON ((0 0, 0 0, 0 0, 0 0, 0 0))',4326),
    provider_id UUID NOT NULL REFERENCES provider(id),
	deleted boolean NOT NULL DEFAULT false
);

-- watershed_usgs_sites
CREATE TABLE IF NOT EXISTS watershed_usgs_sites (
    watershed_id UUID NOT NULL REFERENCES watershed(id),
    usgs_site_parameter_id UUID NOT NULL REFERENCES usgs_site_parameters(id),
    CONSTRAINT watershed_unique_site_param UNIQUE(watershed_id, usgs_site_parameter_id)
);
-- Add comment to describe table
COMMENT ON TABLE watershed_usgs_sites IS 'This is a bridge table.  Each entry represent a watershed/site/parameter that has been requested for USGS data acquisition.';


-------------------
-- SHAPEFILE UPLOAD
-------------------

-- upload_status definition
CREATE TABLE IF NOT EXISTS upload_status (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL
);

-- watershed_shapefile_uploads definition
CREATE TABLE IF NOT EXISTS watershed_shapefile_uploads (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	watershed_id UUID NOT NULL REFERENCES watershed(id),
	file VARCHAR NOT NULL,
	date_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	file_size INTEGER NOT NULL,
    processing_info VARCHAR,
    user_id UUID,
    upload_status_id UUID NOT NULL DEFAULT 'b5d777fc-c46b-4a10-a488-1415e3d7849d' REFERENCES upload_status(id)
);

--------------
-- TIMESERIES
--------------

-- timeseries table
CREATE TABLE IF NOT EXISTS timeseries (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    datasource_id UUID NOT NULL REFERENCES datasource(id),
    datasource_key VARCHAR NOT NULL,
    latest_value DOUBLE PRECISION,
    latest_time TIMESTAMPTZ,
    CONSTRAINT timeseries_unique_datasource UNIQUE(datasource_id, datasource_key)
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
    provider_id UUID NOT NULL REFERENCES provider(id)
);

-- chart_variable_mapping
CREATE TABLE IF NOT EXISTS chart_variable_mapping (
    chart_id UUID NOT NULL REFERENCES chart(id),
    variable VARCHAR NOT NULL,
    timeseries_id UUID NOT NULL REFERENCES timeseries(id) ON DELETE CASCADE,
    CONSTRAINT chart_id_unique_variable UNIQUE(chart_id, variable)
);