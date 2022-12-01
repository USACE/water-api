
-- Enable data acquisition for Watershed/Site/Parameter

INSERT into watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
-- LRH - GUYANDOTTE Watershed
-- LRH - GUYANDOTTE RIVER AT LOGAN, WV - Stage and Flow
('50372dbc-f254-4584-8345-1c3613d2a102', (SELECT usp.id FROM usgs_site_parameters usp JOIN usgs_site us ON us.location_id = usp.location_id JOIN usgs_parameter up ON up.id = usp.parameter_id WHERE us.site_number = '03203600' AND up.code = '00060')),
('50372dbc-f254-4584-8345-1c3613d2a102', (SELECT usp.id FROM usgs_site_parameters usp JOIN usgs_site us ON us.location_id = usp.location_id JOIN usgs_parameter up ON up.id = usp.parameter_id WHERE us.site_number = '03203600' AND up.code = '00065')),
-- LRH - GUYANDOTTE RIVER AT BRANCHLAND, WV - Stage
('50372dbc-f254-4584-8345-1c3613d2a102', (SELECT usp.id FROM usgs_site_parameters usp JOIN usgs_site us ON us.location_id = usp.location_id JOIN usgs_parameter up ON up.id = usp.parameter_id WHERE us.site_number = '03207020' AND up.code = '00065')),


-- LRH - Kanawha Watershed
-- LRH - PINEY CREEK AT RALEIGH, WV (stage, flow)
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (SELECT usp.id FROM usgs_site_parameters usp JOIN usgs_site us ON us.location_id = usp.location_id JOIN usgs_parameter up ON up.id = usp.parameter_id WHERE us.site_number = '03185000' AND up.code = '00060')),
('65a93467-c9b4-4166-acb6-58e8ec06ed3b', (SELECT usp.id FROM usgs_site_parameters usp JOIN usgs_site us ON us.location_id = usp.location_id JOIN usgs_parameter up ON up.id = usp.parameter_id WHERE us.site_number = '03185000' AND up.code = '00065'));