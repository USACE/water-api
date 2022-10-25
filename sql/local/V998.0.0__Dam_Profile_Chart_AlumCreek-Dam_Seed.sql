-- DO $$

-- DECLARE 
-- lrh_cwms_timeseries_id uuid;
-- lrh_cwms_level_id uuid;

-- BEGIN

--     SELECT d.id into lrh_cwms_timeseries_id
--     FROM datasource d 
--     JOIN datasource_type dt ON dt.id = d.datasource_type_id 
--     JOIN provider p ON p.id = d.provider_id 
--     WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'lrh';

--     SELECT d.id into lrh_cwms_level_id
--     FROM datasource d 
--     JOIN datasource_type dt ON dt.id = d.datasource_type_id 
--     JOIN provider p ON p.id = d.provider_id 
--     WHERE dt.slug = 'cwms-level' AND p.slug = 'lrh';

--     -- Alum Creek Dam

--     INSERT into timeseries(datasource_id, datasource_key, latest_time, latest_value) VALUES
--     (lrh_cwms_timeseries_id, 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS', '2022-09-27T17:00:00Z', 888.14),
--     (lrh_cwms_timeseries_id, 'AlumCr-Lake.Flow.Inst.15Minutes.0.OBS', '2022-09-27T16:00:00Z', 23.56),
--     (lrh_cwms_timeseries_id, 'AlumCr-Outflow.Flow.Inst.15Minutes.0.OBS', '2022-09-27T16:00:00Z', 16.56),
--     (lrh_cwms_timeseries_id, 'AlumCr-Outflow.Stage.Inst.15Minutes.0.OBS', '2022-09-27T16:00:00Z', 1.38),
--     (lrh_cwms_level_id, 'AlumCr.Elev.Inst.0.Streambed', '2022-09-09T18:35:00Z', 820),
--     (lrh_cwms_level_id, 'AlumCr.Elev.Inst.0.Top of Dam', '2022-09-09T18:35:00Z', 913),
--     (lrh_cwms_level_id, 'AlumCr.Elev.Inst.0.Top of Flood', '2022-09-09T18:35:00Z', 901);

--     -- Alum Creek Locations already in seed data and should be in develop/stable
--     INSERT INTO chart(id, location_id, slug, name, type_id, provider_id) VALUES
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', (SELECT id from location where slug = 'alumcr'), 'alum-creek-dam-profile', 'Alum Creek Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c', '2f160ba7-fd5f-4716-8ced-4a29f75065a6');

--     INSERT INTO chart_variable_mapping(chart_id, variable, timeseries_id) VALUES
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'pool', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS')),
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'tail', (SELECT id from timeseries where datasource_key = 'AlumCr-Outflow.Stage.Inst.15Minutes.0.OBS')),
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'inflow', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Flow.Inst.15Minutes.0.OBS')),
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'outflow', (SELECT id from timeseries where datasource_key = 'AlumCr-Outflow.Flow.Inst.15Minutes.0.OBS')),
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'streambed', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Streambed')),
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Top of Dam')),
--     ('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Top of Flood'));

-- END$$;