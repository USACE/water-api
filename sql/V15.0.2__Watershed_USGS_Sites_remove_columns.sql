-- Transfer the seed data for watershed_usgs_sites to usgs_site_parameters
-- before consolidation

-- LRH - BLUESTONE LAKE NEAR HINTON, WV (elev NGVD29)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number='03179800'), 'f739b4af-1c96-437c-a788-901f59d177fb')
ON CONFLICT DO NOTHING;
-- LRH - NEW RIVER AT HINTON, WV
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number='03184500'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number='03184500'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
((select id from usgs_site where site_number='03184500'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e') 
ON CONFLICT DO NOTHING;
-- LRH - PINEY CREEK AT RALEIGH, WV (stage, flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number='03185000'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number='03185000'), 'ba29fc34-6315-4424-838f-9b1863715fad') 
ON CONFLICT DO NOTHING;
-- LRH - NEW RIVER AT THURMOND, WV (stage, flow, water temp, precip)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03185400'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number = '03185400'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
((select id from usgs_site where site_number = '03185400'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e'),
((select id from usgs_site where site_number = '03185400'), '738eb4df-b34b-45cc-a5aa-f2136384244f')
ON CONFLICT DO NOTHING;
-- LRH - NEW RIVER BELOW HAWKS NEST DAM, WV (stage)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '380649081083301'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28')
ON CONFLICT DO NOTHING;
-- LRH - KANAWHA RIVER AT KANAWHA FALLS, WV (stage, flow, water temp)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03193000'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number = '03193000'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
((select id from usgs_site where site_number = '03193000'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e')
ON CONFLICT DO NOTHING;
-- LRN - CUMBERLAND RIVER AT NASHVILLE, TN (stage, flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '03431500'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number = '03431500'), 'ba29fc34-6315-4424-838f-9b1863715fad')
ON CONFLICT DO NOTHING;
-- MVP - MISSISSIPPI RIVER AT ST. PAUL, MN (stage, flow, water temp)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '05331000'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number = '05331000'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
((select id from usgs_site where site_number = '05331000'), '0fa9993d-6674-4ba3-ac8a-f02830beea1e')
ON CONFLICT DO NOTHING;
-- MVP - MISSISSIPPI RIVER BELOW L&D #2 AT HASTINGS, MN (stage, flow)
INSERT INTO usgs_site_parameters (site_id, parameter_id) VALUES
((select id from usgs_site where site_number = '05331580'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
((select id from usgs_site where site_number = '05331580'), 'ba29fc34-6315-4424-838f-9b1863715fad')
ON CONFLICT DO NOTHING;

-- Remove usgs_site_id and usgs_parameter_id in table watershed_usgs_sites
-- Replace with id from usgs_site_parameters

-- 1) Add new Column usgs_site_parameter_id
ALTER TABLE watershed_usgs_sites ADD COLUMN "usgs_site_parameter_id" UUID REFERENCES usgs_site_parameters(id);

-- 2) Lookup the usgs_site_parameter_id from usgs_site_parameter that represents the existing site_id and parameter_id
--    Set the value of usgs_site_parameter_id
UPDATE watershed_usgs_sites SET usgs_site_parameter_id = (select id from usgs_site_parameters usp 
   where usp.site_id = usgs_site_id AND usp.parameter_id = usgs_parameter_id);

-- 3) Remove any records where usgs_site_parameter_id lookup resulted in NULL (not found)
DELETE from watershed_usgs_sites where usgs_site_parameter_id is NULL;

-- 4) Remove contraint watershed_unique_site_param
ALTER TABLE watershed_usgs_sites DROP CONSTRAINT watershed_unique_site_param;

-- 5 ) Remove columns usgs_site_id and usgs_parameter_id
ALTER TABLE watershed_usgs_sites DROP COLUMN usgs_site_id;
ALTER TABLE watershed_usgs_sites DROP COLUMN usgs_parameter_id;

-- 6) Add constraint
ALTER TABLE watershed_usgs_sites ADD CONSTRAINT watershed_unique_site_param UNIQUE(watershed_id, usgs_site_parameter_id);

-- Set usgs_site_parameter_id to NOT NULL
ALTER TABLE watershed_usgs_sites ALTER COLUMN usgs_site_parameter_id SET NOT NULL;

-- Add comment to describe table
COMMENT ON TABLE watershed_usgs_sites IS 'This is a bridge table.  Each entry represent a watershed/site/parameter that has been requested for USGS data acquisition.';