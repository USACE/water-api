--additional seed data for testing

INSERT into watershed_usgs_sites (watershed_id, usgs_site_id, usgs_parameter_id) VALUES
-- LRH - BLUESTONE LAKE NEAR HINTON, WV (elev NGVD29)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03179800'), 'f739b4af-1c96-437c-a788-901f59d177fb'),
-- LRH - NEW RIVER AT HINTON, WV (stage, flow, water temp)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03184500'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03184500'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03184500'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e'),
-- LRH - PINEY CREEK AT RALEIGH, WV (stage, flow)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03185000'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03185000'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
-- LRH - NEW RIVER AT THURMOND, WV (stage, flow, water temp, precip)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03185400'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03185400'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03185400'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03185400'), '738eb4df-b34b-45cc-a5aa-f2136384244f'),
-- lRH - NEW RIVER BELOW HAWKS NEST DAM, WV (stage)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '380649081083301'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
-- LRH - KANAWHA RIVER AT KANAWHA FALLS, WV (stage, flow, water temp)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03193000'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03193000'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (select id from usgs_site where site_number = '03193000'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e'),
-- LRN - CUMBERLAND RIVER AT NASHVILLE, TN (stage, flow)
('c785f4de-ab17-444b-b6e6-6f1ad16676e8', (select id from usgs_site where site_number = '03431500'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('c785f4de-ab17-444b-b6e6-6f1ad16676e8', (select id from usgs_site where site_number = '03431500'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
-- MVP - MISSISSIPPI RIVER AT ST. PAUL, MN (stage, flow, water temp)
('048ce853-6642-4ac4-9fb2-81c01f67a85b', (select id from usgs_site where site_number = '05331000'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('048ce853-6642-4ac4-9fb2-81c01f67a85b', (select id from usgs_site where site_number = '05331000'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
('048ce853-6642-4ac4-9fb2-81c01f67a85b', (select id from usgs_site where site_number = '05331000'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e'),
-- MVP - MISSISSIPPI RIVER BELOW L&D #2 AT HASTINGS, MN (stage, flow)
('048ce853-6642-4ac4-9fb2-81c01f67a85b', (select id from usgs_site where site_number = '05331580'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('048ce853-6642-4ac4-9fb2-81c01f67a85b', (select id from usgs_site where site_number = '05331580'), 'ba29fc34-6315-4424-838f-9b1863715fad');

