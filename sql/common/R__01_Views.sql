--------------
-- V_USGS_SITE
--------------

CREATE OR REPLACE VIEW v_usgs_site AS (
     SELECT
        l.id,
        l.code AS site_number,
        l."attributes"->>'station_name' AS "name",
        ST_AsGeoJSON(l.geometry)::json AS geometry,
        sa.stusps AS state_abbrev,
        sa.name as state,
        COALESCE(code_agg.parameter_codes, '{}') as parameter_codes,
        l.create_date,
        l.update_date
        FROM "location" l 
        JOIN datasource d ON d.id = l.datasource_id 
        JOIN datasource_type dt ON dt.id = d.datasource_type_id
        JOIN tiger_data.state_all sa ON sa.gid = l.state_id
        LEFT JOIN (
            SELECT array_agg(code) AS parameter_codes, b.location_id 
            FROM usgs_parameter a
            JOIN usgs_site_parameters b ON b.parameter_id = a.id
            GROUP BY b.location_id
            ) code_agg ON code_agg.location_id = l.id
        WHERE dt.slug = 'usgs-site'
);

--------------
-- V_WATERSHED
--------------
-- This should rebuild after being deleted in 1.0.7
CREATE OR REPLACE VIEW v_watershed AS (
    SELECT w.id,
           w.slug,
           w.name,
           w.geometry AS geometry,
           w.provider_id,
           p.slug AS provider_slug
	FROM   watershed w
    LEFT JOIN provider p ON w.provider_id = p.id
	WHERE NOT w.deleted
);