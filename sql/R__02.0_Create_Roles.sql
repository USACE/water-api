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
ALTER ROLE water_user SET search_path TO a2w_cwms,topology,tiger,tiger_data,public;

-- Add State Column to Location
ALTER TABLE a2w_cwms.location ADD COLUMN state_id integer REFERENCES tiger_data.state_all(gid);