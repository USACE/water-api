-- Shared Function to Notify Cumulus Async Listener Functions (ALF) Listener
CREATE OR REPLACE FUNCTION notify_async_listener(t text) RETURNS void AS $$
    BEGIN
        PERFORM (SELECT pg_notify('water_new', t));
    END;
$$ LANGUAGE plpgsql;

--------------------------------------------------------------
-- ASYNC LISTENER FUNCTION (ALF) FOR shapefile_geoprocess
--------------------------------------------------------------

-- Trigger Function; Inserts Into acquirablefile Table
-- CREATE OR REPLACE FUNCTION notify_shapefile_geoprocess() RETURNS trigger AS $$
--     BEGIN
--         PERFORM (
--             WITH geoprocess_config as (
--                 	SELECT id AS shapefile_upload_id,
--                 	watershed_id,
--                 	(SELECT config_value from config where config_name = 'write_to_bucket') AS bucket,
--                 	file      AS key
--                 FROM watershed_shapefile_uploads
--                 WHERE id = NEW.id
--             )
--             SELECT notify_async_listener(
--                 json_build_object(
--                     'fn', 'geoprocess-shapefile-upload',
--                     'details', json_build_object(
--                         'processor', 'watershed_shapefile_upload',
--                         'input', row_to_json(geoprocess_config),
--                         'functions', array[
--                             json_build_object('function', 'cleanup'), 
--                             json_build_object('function', 'dissolve'),
--                             json_build_object('function', 'simplify'),
--                             json_build_object('function', 'transform')
--                             ],
--                         'output_target', 'watersheds/'||geoprocess_config.watershed_id||'/update_geometry'
--                     )::text
--                 )::text
--             ) FROM geoprocess_config
--         );
--         RETURN NULL;
--     END;
-- $$ LANGUAGE plpgsql;

-- Trigger
-- CREATE TRIGGER notify_shapefile_geoprocess
-- AFTER INSERT ON watershed_shapefile_uploads
-- FOR EACH ROW
-- EXECUTE PROCEDURE notify_shapefile_geoprocess();