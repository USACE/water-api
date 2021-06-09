-- For a production-ready deployment scenario, the role 'water_user' with a complicated selected password
-- should already exist, having been created when the database was stood-up.
-- The statement below is used to create database user for developing locally with Docker Compose with a
-- simple password ('water_pass'). https://stackoverflow.com/questions/8092086/create-postgresql-role-user-if-it-doesnt-exist
DO $$
BEGIN
  CREATE USER water_user WITH ENCRYPTED PASSWORD 'water_pass';
  EXCEPTION WHEN DUPLICATE_OBJECT THEN
  RAISE NOTICE 'not creating role water_user -- it already exists';
END
$$;

-- Role water_reader;
DO $$
BEGIN
  CREATE ROLE water_reader;
  EXCEPTION WHEN DUPLICATE_OBJECT THEN
  RAISE NOTICE 'not creating role water_reader -- it already exists';
END
$$;

-- Role water_writer
DO $$
BEGIN
  CREATE ROLE water_writer;
  EXCEPTION WHEN DUPLICATE_OBJECT THEN
  RAISE NOTICE 'not creating role water_writer -- it already exists';
END
$$;

-- Role postgis_reader
DO $$
BEGIN
  CREATE ROLE postgis_reader;
  EXCEPTION WHEN DUPLICATE_OBJECT THEN
  RAISE NOTICE 'not creating role postgis_reader -- it already exists';
END
$$;

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
