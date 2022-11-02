DO $$

DECLARE 
    lrh_cwms_location_datasource_id uuid;

BEGIN

    SELECT d.id into lrh_cwms_location_datasource_id
    FROM datasource d 
    JOIN datatype dt ON dt.id = d.datatype_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-location' AND p.slug = 'lrh';

    INSERT INTO location (datasource_id, slug, code, geometry, state_id, attributes) VALUES
    (lrh_cwms_location_datasource_id,'lrh-test-project','LRHTestProject',ST_GeomFromText('POINT(-0.0 0.0)',4326),1,'{"public_name":"LRH Test Project","kind":"PROJECT"}'),
    -- LRH - Scioto Basin
    (lrh_cwms_location_datasource_id,'alumcr','AlumCr',ST_GeomFromText('POINT(-0.0 0.0)',4326),25,'{"public_name":"Alum Creek Lake","kind":"PROJECT"}'),
    -- LRH - Additional Muskingum Basin Projects
    (lrh_cwms_location_datasource_id,'beachcity','BeachCity',ST_GeomFromText('POINT(-0.0 0.0)',4326),25,'{"public_name":"Beach City Dam","kind":"PROJECT"}'),
    (lrh_cwms_location_datasource_id,'beachcity-lake','BeachCity-Lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25,'{"public_name":"Beach City Lake","kind":"STREAM_LOCATION"}'),
    (lrh_cwms_location_datasource_id,'beachcity-outflow','BeachCity-Outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25,'{"public_name":"Beach City Outflow","kind":"STREAM_LOCATION"}'),
    (lrh_cwms_location_datasource_id,'bolivar','Bolivar',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25,'{"public_name":"Bolivar Dam","kind":"PROJECT"}'),
    (lrh_cwms_location_datasource_id,'bolivar-lake','Bolivar-Lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25,'{"public_name":"Bolivar Pool","kind":"STREAM_LOCATION"}'),
    (lrh_cwms_location_datasource_id,'bolivar-outflow','Bolivar-Outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25,'{"public_name":"Bolivar Outflow","kind":"STREAM_LOCATION"}'),
    -- LRH - Kanawha Basin Projects
    (lrh_cwms_location_datasource_id,'bluestone','Bluestone',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1,'{"public_name":"Bluestone Dam","kind":"PROJECT"}'),
    (lrh_cwms_location_datasource_id,'bluestone-lake','Bluestone-Lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1,'{"public_name":"Bluestone Lake","kind":"STREAM_LOCATION"}'),
    (lrh_cwms_location_datasource_id,'bluestone-outflow','Bluestone-Outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1,'{"public_name":"Bluestone Outflow","kind":"STREAM_LOCATION"}');

END$$;
