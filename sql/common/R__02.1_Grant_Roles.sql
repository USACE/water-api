-- Ensure this runs everytime by using dynamic timestamp
-- ${flyway:timestamp}

-- Grant Schema Usage to water_user
GRANT USAGE ON SCHEMA water TO water_user;

-- Grant 'tiger' Schema Usage to water_user
GRANT USAGE ON SCHEMA tiger TO water_user;
GRANT SELECT ON tiger.state TO water_user;

-- Grant 'tiger_data' Schema Usage to water_user
GRANT USAGE ON SCHEMA tiger_data TO water_user;
GRANT SELECT ON tiger_data.state_all TO water_user;

--------------------------------------------------------------------------
-- NOTE: IF USERS ALREADY EXIST ON DATABASE, JUST RUN FROM THIS POINT DOWN
--------------------------------------------------------------------------

GRANT SELECT ON
    config,
    datasource,
    datatype,
    location,    
    provider,
    timeseries,
    timeseries_value,
    timeseries_group,
    timeseries_group_members,
    chart,
    chart_variable_mapping,
    v_chart,
    v_chart_detail,
    v_datasource,
    v_location,
    v_timeseries,
    v_timeseries_group,
    v_timeseries_group_detail
TO water_reader;

-- Role water_writer
-- Tables specific to water app
GRANT INSERT,UPDATE,DELETE ON
    config,
    datasource,
    datatype,
    location,    
    provider,
    timeseries,
    timeseries_value,
    timeseries_group,
    timeseries_group_members,
    chart,
    chart_variable_mapping
TO water_writer;

-- Role postgis_reader
GRANT SELECT ON geometry_columns TO postgis_reader;
GRANT SELECT ON geography_columns TO postgis_reader;
GRANT SELECT ON spatial_ref_sys TO postgis_reader;

-- Grant Permissions to water_user
GRANT postgis_reader TO water_user;
GRANT water_reader TO water_user;
GRANT water_writer TO water_user;