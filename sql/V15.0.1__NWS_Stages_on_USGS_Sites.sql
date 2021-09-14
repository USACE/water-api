-- add NWS Flood stages to USGS Sites View

DROP VIEW v_usgs_site;

-- REPLACE VIEW v_usgs_site
CREATE OR REPLACE VIEW v_usgs_site AS (
    SELECT 
    s.id,
    s.site_number,
    s.name,
    ST_AsGeoJSON(s.geometry)::json AS geometry,
    s.elevation,
    s.horizontal_datum_id,
    s.vertical_datum_id,
    v.name AS vertical_datum,
    s.huc,
    s.state_abbrev,
    st.name as state,
    COALESCE(code_agg.parameter_codes, '{}') as parameter_codes,
    json_build_object('action_stage', ns.action_stage, 'flood_stage', ns.flood_stage, 
    	'moderate_flood_stage', ns.moderate_flood_stage, 'major_flood_stage', ns.major_flood_stage) AS nws_stages,
    s.create_date,
    s.update_date
    FROM usgs_site s
    JOIN vertical_datum v ON v.id=s.vertical_datum_id
    JOIN tiger_data.state_all st ON st.stusps=s.state_abbrev
    LEFT JOIN nws_stages ns ON ns.usgs_site_number=s.site_number
    LEFT JOIN (
        SELECT array_agg(code) AS parameter_codes, b.site_id 
        FROM usgs_parameter a
        JOIN usgs_site_parameters b ON b.parameter_id = a.id
        GROUP BY b.site_id
        ) code_agg ON code_agg.site_id = s.id
);

-- Grant read
GRANT SELECT ON v_usgs_site TO water_reader;