-- datasource table
CREATE TABLE IF NOT EXISTS datasource (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    uri VARCHAR NOT NULL
);

-- timeseries table
CREATE TABLE IF NOT EXISTS timeseries (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    datasource_id UUID NOT NULL REFERENCES datasource(id),
    datasource_key VARCHAR UNIQUE NOT NULL,
    latest_value DOUBLE PRECISION NOT NULL,
    latest_time TIMESTAMPTZ NOT NULL
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

INSERT into datasource(id, slug, name, uri) VALUES
('a138e363-30ea-4e0d-8d8f-cce03cb8e1d0', 'lrh-ts-cwms-radar', 'CWMS RADAR', 'https://cwms-data.usace.army.mil/cwms-data/timeseries?office=lrh'),
('5bd097f6-0ecb-4b85-bc95-b0ccb945a8f0', 'mvp-ts-cwms-radar', 'CWMS RADAR', 'https://cwms-data.usace.army.mil/cwms-data/timeseries?office=mvp'),
('36dc9f8c-b18b-433c-b919-9c067739b6aa', 'usgs-webservice', 'USGS Webservice', 'https//waterservices.usgs.gov/nwis/iv/'),
('fa89f10d-8111-405b-8b93-4645ce10b4ac', 'sas-file-reader', 'SAS S3 File Consumer', 'S3://mybucket/mypath/to/sas/file');

INSERT into timeseries(id, datasource_id, datasource_key, latest_time, latest_value) VALUES
('4c4d08d9-356e-4a29-9867-d6c48ca1b7ff', 'a138e363-30ea-4e0d-8d8f-cce03cb8e1d0', 'AlumCr-Lake.Elev.Inst.15Minutes.0.OBS', '2022-09-08 20:00:00-00', 888.22);

INSERT INTO visualization(id, location_id, slug, name, type_id) VALUES
('ff03ac20-439c-4504-8c33-8819db2acb23', (SELECT id from location where slug = 'alumcr'), 'lrh-alum-creek-dam-profile', 'Alum Creek Dam Profile Chart', '53da77d0-6550-4f02-abf8-4bcd1a596a7c');

INSERT INTO visualization_variable_mapping(visualization_id, variable, timeseries_id) VALUES
('ff03ac20-439c-4504-8c33-8819db2acb23', 'pool', '4c4d08d9-356e-4a29-9867-d6c48ca1b7ff');

