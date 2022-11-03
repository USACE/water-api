DO $$

DECLARE 
    lrh_cwms_location_datasource_id uuid;
    lrn_cwms_location_datasource_id uuid;
    mvp_cwms_location_datasource_id uuid;

BEGIN

    SELECT id into lrh_cwms_location_datasource_id FROM v_datasource  WHERE datatype = 'cwms-location' AND provider = 'lrh';
    SELECT id into lrn_cwms_location_datasource_id FROM v_datasource  WHERE datatype = 'cwms-location' AND provider = 'lrn';
    SELECT id into mvp_cwms_location_datasource_id FROM v_datasource  WHERE datatype = 'cwms-location' AND provider = 'mvp';

    INSERT INTO location (datasource_id, slug, code, geometry, state_id, attributes) VALUES
    (lrh_cwms_location_datasource_id,'lrh-test-project','LRHTestProject',ST_GeomFromText('POINT(-0.0 0.0)',4326),1,'{"public_name":"LRH Test Project","kind":"PROJECT"}'),
    -- LRH - Scioto Basin
    (lrh_cwms_location_datasource_id,'alumcr','AlumCr',ST_GeomFromText('POINT(-0.0 0.0)',4326),25,'{"public_name":"Alum Creek Lake","kind":"PROJECT"}'),
    (lrh_cwms_location_datasource_id,'alumcr-lake','AlumCr-Lake',ST_GeomFromText('POINT(-0.0 0.0)',4326),25,'{"public_name":"Alum Creek Lake","kind":"STREAM_LOCATION"}'),
    (lrh_cwms_location_datasource_id,'alumcr-outflow','AlumCr-Outflow',ST_GeomFromText('POINT(-0.0 0.0)',4326),25,'{"public_name":"Alum Creek Outflow","kind":"STREAM_LOCATION"}'),
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
    (lrh_cwms_location_datasource_id,'bluestone-outflow','Bluestone-Outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1,'{"public_name":"Bluestone Outflow","kind":"STREAM_LOCATION"}'),
    -- LRN
    (lrn_cwms_location_datasource_id, 'bahk2-barkley', 'BAHK2-BARKLEY', ST_GeomFromText('POINT(-0.0 0.0)', 4326),1,'{"public_name":"BARKLEY","kind":"PROJECT"}'),
    -- MVP
    (mvp_cwms_location_datasource_id, 'baldhill-dam', 'Baldhill_Dam', ST_GeomFromText('POINT(-0.0 0.0)', 4326),1,'{"public_name":"Baldhill Dam","kind":"PROJECT"}');

END$$;
