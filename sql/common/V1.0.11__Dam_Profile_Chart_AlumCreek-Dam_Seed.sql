-- Alum Creek Dam

INSERT into timeseries(datasource_id, datasource_key, latest_time, latest_value) VALUES
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS', '2022-09-27T17:00:00Z', 888.14),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Lake.Flow.Inst.15Minutes.0.OBS', '2022-09-27T16:00:00Z', 23.56),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Outflow.Flow.Inst.15Minutes.0.OBS', '2022-09-27T16:00:00Z', 16.56),
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Outflow.Stage.Inst.15Minutes.0.OBS', '2022-09-27T16:00:00Z', 1.38),
('5bb6d520-5223-4b04-b348-f57268a41c03', 'AlumCr.Elev.Inst.0.Streambed', '2022-09-09T18:35:00Z', 820),
('5bb6d520-5223-4b04-b348-f57268a41c03', 'AlumCr.Elev.Inst.0.Top of Dam', '2022-09-09T18:35:00Z', 913),
('5bb6d520-5223-4b04-b348-f57268a41c03', 'AlumCr.Elev.Inst.0.Top of Flood', '2022-09-09T18:35:00Z', 901);

-- Alum Creek Locations already in seed data and should be in develop/stable
INSERT INTO visualization(id, location_id, slug, name, type_id) VALUES
('ff03ac20-439c-4504-8c33-8819db2acb23', (SELECT id from location where slug = 'alumcr'), 'alum-creek-dam-profile', 'Alum Creek Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c');

INSERT INTO visualization_variable_mapping(visualization_id, variable, timeseries_id) VALUES
('ff03ac20-439c-4504-8c33-8819db2acb23', 'pool', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS')),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'tail', (SELECT id from timeseries where datasource_key = 'AlumCr-Outflow.Stage.Inst.15Minutes.0.OBS')),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'inflow', (SELECT id from timeseries where datasource_key = 'AlumCr-Lake.Flow.Inst.15Minutes.0.OBS')),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'outflow', (SELECT id from timeseries where datasource_key = 'AlumCr-Outflow.Flow.Inst.15Minutes.0.OBS')),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'streambed', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Streambed')),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-dam', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Top of Dam')),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-flood', (SELECT id from timeseries where datasource_key = 'AlumCr.Elev.Inst.0.Top of Flood'));

