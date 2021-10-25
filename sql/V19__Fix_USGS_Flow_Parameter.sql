-- Remove any daily flow average measurements that may be present
DELETE FROM a2w_cwms.usgs_measurements 
WHERE usgs_site_parameters_id IN (
	SELECT DISTINCT um.usgs_site_parameters_id FROM a2w_cwms.usgs_site_parameters usp
	JOIN a2w_cwms.usgs_measurements um ON um.usgs_site_parameters_id = usp.id 
	WHERE usp.parameter_id = 'ba29fc34-6315-4424-838f-9b1863715fad'
);

-- remove daily flow averages param from watershed_usgs_sites 
DELETE FROM a2w_cwms.watershed_usgs_sites 
WHERE usgs_site_parameter_id IN (
	SELECT DISTINCT id FROM a2w_cwms.usgs_site_parameters usp
	WHERE usp.parameter_id = 'ba29fc34-6315-4424-838f-9b1863715fad'
);

-- delete flow averages param as available param for usgs sites
DELETE FROM a2w_cwms.usgs_site_parameters WHERE parameter_id = 'ba29fc34-6315-4424-838f-9b1863715fad';

-- Remove flow daily average USGS parameter
DELETE FROM usgs_parameter where id = 'ba29fc34-6315-4424-838f-9b1863715fad';




-- add flow instantaneous parameter
INSERT INTO usgs_parameter (id, code, description) VALUES 
('06cca640-f52b-4c0c-8f64-a985836fda5a', '00061', 'Discharge, instantaneous, cubic feet per second');

-- re-add inst flow for initial seed data

-- LRH - NEW RIVER AT HINTON, WV (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number='03184500'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- LRH - PINEY CREEK AT RALEIGH, WV (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number='03185000'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- LRH - NEW RIVER AT THURMOND, WV (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03185400'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- LRH - KANAWHA RIVER AT KANAWHA FALLS, WV (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03193000'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- LRH - EAST FORK TWELVEPOLE CREEK NEAR DUNLOW, WV (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03206600'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- LRN - CUMBERLAND RIVER AT NASHVILLE, TN (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03431500'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- MVP - MISSISSIPPI RIVER AT ST. PAUL, MN (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '05331000'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;
-- MVP - MISSISSIPPI RIVER BELOW L&D #2 AT HASTINGS, MN (flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '05331580'), '06cca640-f52b-4c0c-8f64-a985836fda5a')
ON CONFLICT DO NOTHING;

-- re-add watershed configs


-- LRH - NEW RIVER AT HINTON, WV (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '03184500' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
-- LRH - PINEY CREEK AT RALEIGH, WV (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '03185000' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
-- LRH - NEW RIVER AT THURMOND, WV (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '03185400' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
-- LRH - KANAWHA RIVER AT KANAWHA FALLS, WV (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '03193000' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
-- LRH - EAST FORK TWELVEPOLE CREEK NEAR DUNLOW, WV (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('4d3083d1-101c-4b76-9311-1154917ffbf1', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '03206600' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
--
-- LRN - CUMBERLAND RIVER AT NASHVILLE, TN (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('c785f4de-ab17-444b-b6e6-6f1ad16676e8', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '03431500' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
--
-- MVP - MISSISSIPPI RIVER AT ST. PAUL, MN (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('048ce853-6642-4ac4-9fb2-81c01f67a85b', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '05331000' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;
-- MVP - MISSISSIPPI RIVER BELOW L&D #2 AT HASTINGS, MN (flow)
INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
('048ce853-6642-4ac4-9fb2-81c01f67a85b', 
	(SELECT usp.id FROM a2w_cwms.usgs_site us 
	JOIN a2w_cwms.usgs_site_parameters usp ON usp.site_id = us.id
	WHERE us.site_number = '05331580' AND usp.parameter_id = '06cca640-f52b-4c0c-8f64-a985836fda5a')
) ON CONFLICT DO NOTHING;