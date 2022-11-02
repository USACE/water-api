DO $$

DECLARE 
    nws_level_ds_id uuid;

BEGIN

    SELECT d.id into nws_level_ds_id
    FROM datasource d 
    JOIN datatype dt ON dt.id = d.datatype_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'nws-level' AND p.slug = 'noaa-nws';


    -- Kanawha Falls at Kanawha River, WV
    INSERT into timeseries(datasource_id, datasource_key, latest_time, latest_value) VALUES
        (nws_level_ds_id, 'KANW2.Stage.Inst.0.Action', '2021-09-19T18:35:00Z', 25),
        (nws_level_ds_id, 'KANW2.Stage.Inst.0.Flood', '2021-09-19T18:35:00Z', 27),
        (nws_level_ds_id, 'KANW2.Stage.Inst.0.Moderate Flood', '2021-09-19T18:35:00Z', 30),
        (nws_level_ds_id, 'KANW2.Stage.Inst.0.Major Flood', '2021-09-19T18:35:00Z', 35);

END$$;