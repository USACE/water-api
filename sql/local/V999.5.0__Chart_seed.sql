
DO $$

DECLARE 
lrh_cwms_timeseries_ds_id uuid;
lrh_cwms_level_ds_id uuid;

BEGIN

    -- SELECT d.id into lrh_cwms_timeseries_ds_id
    -- FROM datasource d 
    -- JOIN datatype dt ON dt.id = d.datatype_id 
    -- JOIN provider p ON p.id = d.provider_id 
    -- WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'lrh';

    -- SELECT d.id into lrh_cwms_level_ds_id
    -- FROM datasource d 
    -- JOIN datatype dt ON dt.id = d.datatype_id 
    -- JOIN provider p ON p.id = d.provider_id 
    -- WHERE dt.slug = 'cwms-level' AND p.slug = 'lrh';

    -- INSERT into timeseries(datasource_id, datasource_key, latest_time, latest_value) VALUES
    -- (lrh_cwms_timeseries_ds_id, 'Atwood-Lake.Elev.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Atwood-Lake.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Atwood-Outflow.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Atwood-Outflow.Stage.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_level_ds_id, 'Atwood.Elev.Inst.0.Streambed', '2022-09-16T18:35:00Z', 886),
    -- (lrh_cwms_level_ds_id, 'Atwood.Elev.Inst.0.Top of Dam', '2022-09-16T18:35:00Z', 955),
    -- (lrh_cwms_level_ds_id, 'Atwood.Elev.Inst.0.Top of Flood', '2022-09-16T18:35:00Z', 941),
    -- (lrh_cwms_timeseries_ds_id, 'BeachCity-Lake.Elev.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'BeachCity-Lake.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'BeachCity-Outflow.Stage.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'BeachCity-Outflow.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_level_ds_id, 'BeachCity.Elev.Inst.0.Streambed', '2022-09-19T18:35:00Z', 931),
    -- (lrh_cwms_level_ds_id, 'BeachCity.Elev.Inst.0.Top of Dam', '2022-09-19T18:35:00Z', 999.7),
    -- (lrh_cwms_level_ds_id, 'BeachCity.Elev.Inst.0.Top of Flood', '2022-09-19T18:35:00Z', 976.5),
    -- (lrh_cwms_timeseries_ds_id, 'Bolivar-Lake.Elev.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Bolivar-Lake.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Bolivar-Outflow.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_level_ds_id, 'Bolivar.Elev.Inst.0.Streambed', '2022-09-19T18:35:00Z', 890),
    -- (lrh_cwms_level_ds_id, 'Bolivar.Elev.Inst.0.Top of Dam', '2022-09-19T18:35:00Z', 985.5),
    -- (lrh_cwms_level_ds_id, 'Bolivar.Elev.Inst.0.Top of Flood', '2022-09-19T18:35:00Z', 962),
    -- (lrh_cwms_timeseries_ds_id, 'Bluestone-Lake.Elev.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Bluestone-Lake.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Bluestone-Outflow.Flow.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_timeseries_ds_id, 'Bluestone-Outflow.Stage.Inst.15Minutes.0.OBS', null, null),
    -- (lrh_cwms_level_ds_id, 'Bluestone.Elev.Inst.0.Streambed', '2022-09-18T18:35:00Z', 1368),
    -- (lrh_cwms_level_ds_id, 'Bluestone.Elev.Inst.0.Top of Dam', '2022-09-18T18:35:00Z', 1535),
    -- (lrh_cwms_level_ds_id, 'Bluestone.Elev.Inst.0.Top of Flood', '2022-09-18T18:35:00Z', 1520);

    INSERT INTO chart(id, location_id, slug, name, type_id, provider_id) VALUES
    -- alum-creek-dam-profile already in seed data
    ('ff03ac20-439c-4504-8c33-8819db2acb23', (SELECT id from v_location WHERE lower(code) = 'alumcr' and datatype = 'cwms-location' AND provider = 'lrh'), 'alum-creek-dam-profile', 'Alum Creek Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c', '2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('1c89c8f2-b7ef-4cc6-a877-ac6c21d48e87', (SELECT id from v_location WHERE lower(code) = 'alumcr' and datatype = 'cwms-location' AND provider = 'lrh'), 'lrh-alum-creek-example-scatter', 'Alum Creek Dam Scatter Chart', '61910b8c-4dfb-4343-affb-d478b6bf915f', '2f160ba7-fd5f-4716-8ced-4a29f75065a6');
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', (SELECT id from location where slug = 'atwood-1'), 'lrh-atwood-dam-profile', 'Atwood Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c', '2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', (SELECT id from location where slug = 'beachcity'), 'lrh-beachcity-dam-profile', 'Beach City Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c', '2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', (SELECT id from location where slug = 'bolivar'), 'lrh-bolivar-dam-profile', 'Bolivar Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c', '2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', (SELECT id from location where slug = 'bluestone'), 'lrh-bluestone-dam-profile', 'Bluestone Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c', '2f160ba7-fd5f-4716-8ced-4a29f75065a6');

    INSERT INTO chart_variable_mapping(chart_id, variable, timeseries_id) VALUES
    -- alumcr dam profile mapping already in seed data
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'pool', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS')),
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'tail', (SELECT id from timeseries where datasource_key = 'AlumCr-Outflow.Stage.Inst.15Minutes.0.OBS')),
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'inflow', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Flow.Inst.15Minutes.0.OBS')),
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'outflow', (SELECT id from timeseries where datasource_key = 'AlumCr-Outflow.Flow.Inst.15Minutes.0.OBS')),
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'streambed', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Streambed')),
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Top of Dam')),
    ('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Top of Flood'));
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', 'pool', (SELECT id from timeseries where datasource_key = 'Atwood-Lake.Elev.Inst.15Minutes.0.OBS')),
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', 'inflow', (SELECT id from timeseries where datasource_key = 'Atwood-Lake.Flow.Inst.15Minutes.0.OBS')),
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', 'outflow', (SELECT id from timeseries where datasource_key = 'Atwood-Outflow.Flow.Inst.15Minutes.0.OBS')),
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', 'streambed', (SELECT id from timeseries where datasource_key = 'Atwood.Elev.Inst.0.Streambed')),
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'Atwood.Elev.Inst.0.Top of Dam')),
    -- ('41f2f75e-9472-4339-bac6-be898e809aee', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'Atwood.Elev.Inst.0.Top of Flood')),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', 'pool', (SELECT id from timeseries where datasource_key = 'BeachCity-Lake.Elev.Inst.15Minutes.0.OBS')),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', 'inflow', (SELECT id from timeseries where datasource_key = 'BeachCity-Lake.Flow.Inst.15Minutes.0.OBS')),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', 'outflow', (SELECT id from timeseries where datasource_key = 'BeachCity-Outflow.Flow.Inst.15Minutes.0.OBS')),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', 'streambed', (SELECT id from timeseries where datasource_key = 'BeachCity.Elev.Inst.0.Streambed')),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'BeachCity.Elev.Inst.0.Top of Dam')),
    -- ('ec15892b-bfde-421f-adfa-9f7c0f906ec2', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'BeachCity.Elev.Inst.0.Top of Flood')),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', 'pool', (SELECT id from timeseries where datasource_key = 'Bolivar-Lake.Elev.Inst.15Minutes.0.OBS')),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', 'inflow', (SELECT id from timeseries where datasource_key = 'Bolivar-Lake.Flow.Inst.15Minutes.0.OBS')),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', 'outflow', (SELECT id from timeseries where datasource_key = 'Bolivar-Outflow.Flow.Inst.15Minutes.0.OBS')),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', 'streambed', (SELECT id from timeseries where datasource_key = 'Bolivar.Elev.Inst.0.Streambed')),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'Bolivar.Elev.Inst.0.Top of Dam')),
    -- ('8cc22bfd-678c-48ad-a35a-ff427c481cc6', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'Bolivar.Elev.Inst.0.Top of Flood')),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', 'pool', (SELECT id from timeseries where datasource_key = 'Bluestone-Lake.Elev.Inst.15Minutes.0.OBS')),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', 'inflow', (SELECT id from timeseries where datasource_key = 'Bluestone-Lake.Flow.Inst.15Minutes.0.OBS')),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', 'outflow', (SELECT id from timeseries where datasource_key = 'Bluestone-Outflow.Flow.Inst.15Minutes.0.OBS')),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', 'streambed', (SELECT id from timeseries where datasource_key = 'Bluestone.Elev.Inst.0.Streambed')),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'Bluestone.Elev.Inst.0.Top of Dam')),
    -- ('945cd8de-39ec-4a75-836d-f14c0f609571', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'Bluestone.Elev.Inst.0.Top of Flood')),
    -- -- scatter seed
    -- ('1c89c8f2-b7ef-4cc6-a877-ac6c21d48e87', 'pool', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS'));

END$$;