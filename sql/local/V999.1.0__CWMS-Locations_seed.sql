DO $$

DECLARE 
    lrh_cwms_location_datasource_id uuid;

BEGIN

    SELECT d.id into lrh_cwms_location_datasource_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-location' AND p.slug = 'lrh';

    INSERT INTO location (datasource_id, slug, geometry, state_id) VALUES
    (lrh_cwms_location_datasource_id,'lrh-test-project',ST_GeomFromText('POINT(-0.0 0.0)',4326),1),
    -- LRH - Scioto Basin
    (lrh_cwms_location_datasource_id,'alumcr',ST_GeomFromText('POINT(-0.0 0.0)',4326),25),
    -- LRH - Additional Muskingum Basin Projects
    (lrh_cwms_location_datasource_id,'beachcity',ST_GeomFromText('POINT(-0.0 0.0)',4326),25),
    (lrh_cwms_location_datasource_id,'beachcity-lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
    (lrh_cwms_location_datasource_id,'beachcity-outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
    (lrh_cwms_location_datasource_id,'bolivar',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
    (lrh_cwms_location_datasource_id,'bolivar-lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
    (lrh_cwms_location_datasource_id,'bolivar-outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
    -- LRH - Kanawha Basin Projects
    (lrh_cwms_location_datasource_id,'bluestone',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1),
    (lrh_cwms_location_datasource_id,'bluestone-lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1),
    (lrh_cwms_location_datasource_id,'bluestone-outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1);

    INSERT INTO cwms_location (location_id, name, public_name, kind_id) VALUES
    ((SELECT id FROM location where slug = 'lrh-test-project'), 'LRHTestProject', 'LRH Test Project', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'alumcr'), 'AlumCr', 'Alum Creek Lake', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'beachcity'), 'BeachCity', 'Beach City Dam', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'beachcity-lake'), 'BeachCity-Lake', 'Beach City Lake', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'beachcity-outflow'), 'BeachCity-Outflow', 'Beach City Outflow', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bolivar'), 'Bolivar', 'Bolivar Dam', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'bolivar-lake'), 'Bolivar-Lake', 'Bolivar Pool', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bolivar-outflow'), 'Bolivar-Outflow', 'Bolivar Outflow', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bluestone'), 'Bluestone', 'Bluestone Dam', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'bluestone-lake'), 'Bluestone-Lake', 'Bluestone Lake', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bluestone-outflow'), 'Bluestone-Outflow', 'Bluestone Outflow', '1e77acaf-fdee-4e7c-b659-101bce76a229');

END$$;
