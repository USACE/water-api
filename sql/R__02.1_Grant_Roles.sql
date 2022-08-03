-- Ensure this runs everytime by using dynamic timestamp
-- ${flyway:timestamp}

-- Grant Schema Usage to water_user
GRANT USAGE ON SCHEMA a2w_cwms TO water_user;

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
    location,
    location_kind,
    level,
    level_kind,
    level_value,
    nws_stages,
    office,
    upload_status,  
    usgs_site,
    usgs_huc2,
    usgs_huc2_simple,
    usgs_huc4,
    usgs_huc4_simple,
    usgs_huc6,
    usgs_huc6_simple,
    usgs_huc8,
    usgs_huc8_simple,
    usgs_measurements,
    usgs_site_parameters,
    usgs_parameter,
    vertical_datum,
    watershed,
    watershed_shapefile_uploads,
    watershed_usgs_sites,
    v_usgs_site,
    v_watershed
TO water_reader;

-- Role cumulus_writer
-- Tables specific to instrumentation app
GRANT INSERT,UPDATE,DELETE ON
    config,
    location,
    location_kind,
    level,
    level_kind,
    level_value,
    vertical_datum,
    nws_stages,
    upload_status,
    usgs_site,
    usgs_measurements,
    usgs_site_parameters,
    usgs_parameter,
    watershed,
    watershed_shapefile_uploads,
    watershed_usgs_sites
TO water_writer;

-- Role postgis_reader
GRANT SELECT ON geometry_columns TO postgis_reader;
GRANT SELECT ON geography_columns TO postgis_reader;
GRANT SELECT ON spatial_ref_sys TO postgis_reader;

-- Grant Permissions to instrument_user
GRANT postgis_reader TO water_user;
GRANT water_reader TO water_user;
GRANT water_writer TO water_user;