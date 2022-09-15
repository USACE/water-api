-- datasource_type table
CREATE TABLE IF NOT EXISTS datasource_type (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    uri VARCHAR NOT NULL,
    CONSTRAINT datasource_type_slug_unique_uri UNIQUE(slug, uri)
);

-- datasource table
CREATE TABLE IF NOT EXISTS datasource (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES provider(id),
    datasource_type_id UUID NOT NULL REFERENCES datasource_type(id),
    CONSTRAINT datasource_unique_provider_id UNIQUE(provider_id, datasource_type_id)
);

-- timeseries table
CREATE TABLE IF NOT EXISTS timeseries (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    datasource_id UUID NOT NULL REFERENCES datasource(id),
    datasource_key VARCHAR NOT NULL,
    latest_value DOUBLE PRECISION,
    latest_time TIMESTAMPTZ,
    CONSTRAINT timeseries_unique_datasource UNIQUE(datasource_id, datasource_key)
);

-- visualization table
CREATE TABLE IF NOT EXISTS visualization (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    location_id UUID NOT NULL REFERENCES location(id),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    type_id UUID NOT NULL,
    CONSTRAINT visualization_location_unique_slug UNIQUE(location_id, slug)
);

-- visualization_variable_mapping
CREATE TABLE IF NOT EXISTS visualization_variable_mapping (
    visualization_id UUID NOT NULL REFERENCES visualization(id),
    variable VARCHAR NOT NULL,
    timeseries_id UUID NOT NULL REFERENCES timeseries(id),
    CONSTRAINT visualization_id_unique_variable UNIQUE(visualization_id, variable)
);

INSERT into datasource_type(id, slug, name, uri) VALUES
('a138e363-30ea-4e0d-8d8f-cce03cb8e1d0', 'cwms-timeseries', 'CWMS RADAR Timeseries', 'https://cwms-data.usace.army.mil/cwms-data/timeseries'),
('97920d27-ee54-4d35-aec4-c01ec31221a2', 'cwms-levels', 'CWMS RADAR Levels', 'https://cwms-data.usace.army.mil/cwms-data/levels'),
('36dc9f8c-b18b-433c-b919-9c067739b6aa', 'usgs-webservice', 'USGS Webservice', 'https//waterservices.usgs.gov/nwis/iv/');

INSERT into datasource(id, provider_id, datasource_type_id) VALUES
('9680cd77-f2fd-47d1-ac29-d71ec4310ea7', '2f160ba7-fd5f-4716-8ced-4a29f75065a6', 'a138e363-30ea-4e0d-8d8f-cce03cb8e1d0'),
('5bb6d520-5223-4b04-b348-f57268a41c03', '2f160ba7-fd5f-4716-8ced-4a29f75065a6', '97920d27-ee54-4d35-aec4-c01ec31221a2'),
('bc8115f1-9570-4153-853b-901518549600', '552e59f7-c0cc-4689-8a4d-e791c028430a', 'a138e363-30ea-4e0d-8d8f-cce03cb8e1d0');

INSERT into timeseries(id, datasource_id, datasource_key, latest_time, latest_value) VALUES
('4c4d08d9-356e-4a29-9867-d6c48ca1b7ff', '9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS', null, null),
('dc158cb4-45b4-42ed-8d8e-9e1cc57aa37c', '9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Lake.Flow.Inst.15Minutes.0.OBS', null, null),
('6af8bdaf-8759-43ef-85ec-070bc8f04220', '9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Outflow.Flow.Inst.15Minutes.0.OBS', null, null),
('2af62690-01b6-4d3a-b234-9fdda9e5e106', '9680cd77-f2fd-47d1-ac29-d71ec4310ea7', 'AlumCr-Outflow.Stage.Inst.15Minutes.0.OBS', null, null),
('810388ce-45e0-4434-86c5-4a1a5883879d', '5bb6d520-5223-4b04-b348-f57268a41c03', 'AlumCr.Elev.Inst.0.Streambed', '2022-09-09T18:35:00-00', 820),
('3a3791ec-b7ca-492d-bac9-6cf4976fd049', '5bb6d520-5223-4b04-b348-f57268a41c03', 'AlumCr.Elev.Inst.0.Top of Dam', '2022-09-09T18:35:00-00', 913),
('672f7a87-bc5e-472b-8ff3-7a1cd9267a03', '5bb6d520-5223-4b04-b348-f57268a41c03', 'AlumCr.Elev.Inst.0.Top of Flood', '2022-09-09T18:35:00-00', 901);

INSERT INTO visualization(id, location_id, slug, name, type_id) VALUES
('ff03ac20-439c-4504-8c33-8819db2acb23', (SELECT id from location where slug = 'alumcr'), 'lrh-alum-creek-dam-profile', 'Alum Creek Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c');

INSERT INTO visualization_variable_mapping(visualization_id, variable, timeseries_id) VALUES
('ff03ac20-439c-4504-8c33-8819db2acb23', 'pool', '4c4d08d9-356e-4a29-9867-d6c48ca1b7ff'),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'inflow', 'dc158cb4-45b4-42ed-8d8e-9e1cc57aa37c'),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'outflow', '6af8bdaf-8759-43ef-85ec-070bc8f04220'),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'streambed', '810388ce-45e0-4434-86c5-4a1a5883879d'),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-dam', '3a3791ec-b7ca-492d-bac9-6cf4976fd049'),
('ff03ac20-439c-4504-8c33-8819db2acb23', 'top-of-flood', '672f7a87-bc5e-472b-8ff3-7a1cd9267a03');

