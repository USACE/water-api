-- level_kind seed data
INSERT INTO level_kind (id, slug, name) VALUES
    ('43e6ecff-32d0-4e03-ba79-f05a9ed5924d', 'top-of-flood-control', 'Top of Flood Control'),
    ('7a998105-2d91-4b2d-ab5d-9d2fe12b9125', 'bottom-of-flood-control', 'Bottom of Flood Control'),
    ('b3e8fbb0-ae51-4f56-b2b0-b39658f72375', 'top-of-hydropower', 'Top of Hydropower');

-- level seed data
INSERT INTO level (id, location_id, level_kind_id) VALUES
    ('bfc583bd-e209-4d9b-8904-c10ce162417e', (SELECT id FROM location WHERE slug = 'dale-hollow'), '43e6ecff-32d0-4e03-ba79-f05a9ed5924d'),
    ('c57c598a-49aa-4321-8dc6-0af79ded687f', (SELECT id FROM location WHERE slug = 'dale-hollow'), '7a998105-2d91-4b2d-ab5d-9d2fe12b9125'),
    ('a0ef34b5-72a9-4ab1-b9d7-ea0fc1b9447f', (SELECT id FROM location WHERE slug = 'dale-hollow'), 'b3e8fbb0-ae51-4f56-b2b0-b39658f72375');

-- level_value seed data for Dale Hollow Lake
-- Julian Date for 1972-01-01 and 1972-12-31
INSERT INTO level_value (level_id, julian_date, value) VALUES
    ('bfc583bd-e209-4d9b-8904-c10ce162417e', (SELECT EXTRACT(EPOCH FROM '1900-01-01 06:00:00-06:00'::timestamptz)), 663), -- Top of Flood Control
    ('bfc583bd-e209-4d9b-8904-c10ce162417e', (SELECT EXTRACT(EPOCH FROM '1900-12-31 06:00:00-06:00'::timestamptz)), 663),
    ('c57c598a-49aa-4321-8dc6-0af79ded687f', (SELECT EXTRACT(EPOCH FROM '1900-01-01 06:00:00-06:00'::timestamptz)), 651), -- Bottom of Flood Control
    ('c57c598a-49aa-4321-8dc6-0af79ded687f', (SELECT EXTRACT(EPOCH FROM '1900-12-31 06:00:00-06:00'::timestamptz)), 651),
    ('a0ef34b5-72a9-4ab1-b9d7-ea0fc1b9447f', (SELECT EXTRACT(EPOCH FROM '1900-01-01 06:00:00-06:00'::timestamptz)), 651), -- Top of Hydropower
    ('a0ef34b5-72a9-4ab1-b9d7-ea0fc1b9447f', (SELECT EXTRACT(EPOCH FROM '1900-12-31 06:00:00-06:00'::timestamptz)), 651);

