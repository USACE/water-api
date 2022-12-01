DO $$

DECLARE 
    nws_level_ds_id uuid;

BEGIN

    SELECT id into nws_level_ds_id 
    FROM v_datasource 
    WHERE datatype = 'nws-level' AND provider = 'noaa-nws';


    -- Kanawha Falls at Kanawha River, WV
    INSERT into timeseries(datasource_id, datasource_key, location_id, latest_time, latest_value) VALUES
        (nws_level_ds_id, 'KANW2.Stage.Action',         (SELECT id from v_location where lower(code) = 'kanw2' and datatype = 'nws-site'), '2021-09-19T18:35:00Z', 25),
        (nws_level_ds_id, 'KANW2.Stage.Flood',          (SELECT id from v_location where lower(code) = 'kanw2' and datatype = 'nws-site'), '2021-09-19T18:35:00Z', 27),
        (nws_level_ds_id, 'KANW2.Stage.Moderate Flood', (SELECT id from v_location where lower(code) = 'kanw2' and datatype = 'nws-site'), '2021-09-19T18:35:00Z', 30),
        (nws_level_ds_id, 'KANW2.Stage.Major Flood',    (SELECT id from v_location where lower(code) = 'kanw2' and datatype = 'nws-site'), '2021-09-19T18:35:00Z', 35);

END$$;