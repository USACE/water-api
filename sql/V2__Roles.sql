CREATE USER water_user WITH ENCRYPTED PASSWORD 'water_pass';
CREATE ROLE water_reader;
CREATE ROLE water_writer;
CREATE ROLE postgis_reader;

-- Set Search Path
ALTER ROLE water_user SET search_path TO a2w_cwms,topology,public;

-- Grant Schema Usage to water_user
GRANT USAGE ON SCHEMA a2w_cwms TO water_user;

--------------------------------------------------------------------------
-- NOTE: IF USERS ALREADY EXIST ON DATABASE, JUST RUN FROM THIS POINT DOWN
--------------------------------------------------------------------------

GRANT SELECT ON
    location,
    location_kind
TO water_reader;

-- Role cumulus_writer
-- Tables specific to instrumentation app
GRANT INSERT,UPDATE,DELETE ON
    location,
    location_kind
TO water_writer;

-- Role postgis_reader
GRANT SELECT ON geometry_columns TO postgis_reader;
GRANT SELECT ON geography_columns TO postgis_reader;
GRANT SELECT ON spatial_ref_sys TO postgis_reader;

-- Grant Permissions to instrument_user
GRANT postgis_reader TO water_user;
GRANT water_reader TO water_user;
GRANT water_writer TO water_user;
