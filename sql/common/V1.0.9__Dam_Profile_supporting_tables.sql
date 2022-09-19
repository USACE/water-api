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

