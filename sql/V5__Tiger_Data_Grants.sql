-- Set Search Path
ALTER ROLE water_user SET search_path TO a2w_cwms,topology,tiger,tiger_data,public;

-- Grant 'tiger' Schema Usage to water_user
GRANT USAGE ON SCHEMA tiger TO water_user;
GRANT SELECT ON tiger.state TO water_user;

-- Grant 'tiger_data' Schema Usage to water_user
GRANT USAGE ON SCHEMA tiger_data TO water_user;
GRANT SELECT ON tiger_data.state_all TO water_user;

-- Add State Column to Location
ALTER TABLE a2w_cwms.location ADD COLUMN state_id integer REFERENCES tiger_data.state_all(gid);
