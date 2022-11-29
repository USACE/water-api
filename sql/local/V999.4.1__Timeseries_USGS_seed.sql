
DO $$

DECLARE 
    -- datasource_id(s)
    usgs_timeseries_ds_id uuid;

    -- location_id(s)
    loc_03203600 uuid;
    loc_03185000 uuid;
    loc_03184500 uuid;

BEGIN
    
    -- Datasource ID
    SELECT id into usgs_timeseries_ds_id FROM v_datasource WHERE datatype = 'usgs-timeseries' AND provider = 'usgs';

    -- Locations
    ------------
    -- GUYANDOTTE RIVER AT LOGAN, WV (03203600)
    SELECT id into loc_03203600 FROM v_location WHERE lower(code) = '03203600' and datatype = 'usgs-site' AND provider = 'usgs';
    -- Piney Creek at Raleigh, WV (03185000)
    SELECT id into loc_03185000 FROM v_location WHERE lower(code) = '03185000' and datatype = 'usgs-site' AND provider = 'usgs';
    -- New River at Hinton, WV (03184500)
    SELECT id into loc_03184500 FROM v_location WHERE lower(code) = '03184500' and datatype = 'usgs-site' AND provider = 'usgs';


    INSERT into timeseries (location_id, datasource_id, datasource_key, latest_time, latest_value) VALUES
    -- GUYANDOTTE RIVER AT LOGAN, WV (03203600)
    (loc_03203600, usgs_timeseries_ds_id, '03203600.00065', '2022-09-27T17:00:00Z', 888.14),
    -- Piney Creek at Raleigh, WV (03185000 00060)
    (loc_03185000, usgs_timeseries_ds_id, '03185000.00060', '2022-09-27T16:00:00Z', 23.56),
    -- Piney Creek at Raleigh, WV (03185000 00065)
    (loc_03185000, usgs_timeseries_ds_id, '03185000.00065', '2022-09-27T16:00:00Z', 23.56),
    -- New River at Hinton, WV (03184500 00065)
    (loc_03184500, usgs_timeseries_ds_id, '03184500.00065', '2022-09-27T16:00:00Z', 23.56);

END$$;
