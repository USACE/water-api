DO $$

DECLARE 

    ts_03203600_00065 uuid; -- GUYANDOTTE RIVER AT LOGAN, WV (03203600) Discharge
    ts_03185000_00060 uuid; -- Piney Creek at Raleigh, WV    (03185000) Gage Height
    ts_03185000_00065 uuid; -- Piney Creek at Raleigh, WV    (03185000) Discharge
    ts_03184500_00065 uuid; -- New River at Hinton, WV       (03184500) Discharge

BEGIN

    SELECT id INTO ts_03203600_00065 FROM v_timeseries WHERE provider = 'usgs' AND datatype = 'usgs-timeseries' AND key = '03203600.00065';
    SELECT id INTO ts_03185000_00060 FROM v_timeseries WHERE provider = 'usgs' AND datatype = 'usgs-timeseries' AND key = '03185000.00060';
    SELECT id INTO ts_03185000_00065 FROM v_timeseries WHERE provider = 'usgs' AND datatype = 'usgs-timeseries' AND key = '03185000.00065';
    SELECT id INTO ts_03184500_00065 FROM v_timeseries WHERE provider = 'usgs' AND datatype = 'usgs-timeseries' AND key = '03184500.00065';
    
    -- SAMPLE TIMESERIES GROUP FOR PROVIDER LRH CONTAINING 4 USGS TIMESERIES THAT HAVE EXAMPLE VALUES
    INSERT INTO timeseries_group (id, slug, name, provider_id) VALUES 
        ('920bdea5-2f3d-41ef-9c59-d5689d3cab30', 'test-timeseries-group', 'Test Timeseries Group', (SELECT id FROM provider WHERE slug = 'lrh'));
    
    INSERT INTO timeseries_group_members (timeseries_group_id, timeseries_id) VALUES
        ('920bdea5-2f3d-41ef-9c59-d5689d3cab30', ts_03203600_00065),
        ('920bdea5-2f3d-41ef-9c59-d5689d3cab30', ts_03185000_00060),
        ('920bdea5-2f3d-41ef-9c59-d5689d3cab30', ts_03185000_00065),
        ('920bdea5-2f3d-41ef-9c59-d5689d3cab30', ts_03184500_00065);

END$$;