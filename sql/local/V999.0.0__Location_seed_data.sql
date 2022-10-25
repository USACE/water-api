INSERT INTO location (datasource_id, slug, geometry, state_id) VALUES
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','lrh-test-project',ST_GeomFromText('POINT(-0.0 0.0)',4326),1),
-- LRH - Additional Muskingum Basin Projects
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','beachcity',ST_GeomFromText('POINT(-0.0 0.0)',4326),25),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','beachcity-lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','beachcity-outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','bolivar',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','bolivar-lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','bolivar-outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),25),
-- LRH - Kanawha Basin Projects
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','bluestone',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','bluestone-lake',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7','bluestone-outflow',ST_GeomFromText('POINT(-0.0 0.0)', 4326),1);

INSERT INTO cwms_location (location_id, name, public_name, kind_id) VALUES
    ((SELECT id FROM location where slug = 'lrh-test-project'), 'LRHTestProject', 'LRH Test Project', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'beachcity'), 'BeachCity', 'Beach City Dam', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'beachcity-lake'), 'BeachCity-Lake', 'Beach City Lake', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'beachcity-outflow'), 'BeachCity-Outflow', 'Beach City Outflow', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bolivar'), 'Bolivar', 'Bolivar Dam', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'bolivar-lake'), 'Bolivar-Lake', 'Bolivar Pool', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bolivar-outflow'), 'Bolivar-Outflow', 'Bolivar Outflow', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bluestone'), 'Bluestone', 'Bluestone Dam', '460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ((SELECT id FROM location where slug = 'bluestone-lake'), 'Bluestone-Lake', 'Bluestone Lake', '1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ((SELECT id FROM location where slug = 'bluestone-outflow'), 'Bluestone-Outflow', 'Bluestone Outflow', '1e77acaf-fdee-4e7c-b659-101bce76a229');

