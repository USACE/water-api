--------------
-- V_USGS_SITE
--------------

-- CREATE OR REPLACE VIEW v_usgs_site AS (
--     SELECT 
--     s.id,
--     s.site_number,
--     s.name,
--     ST_AsGeoJSON(s.geometry)::json AS geometry,
--     s.elevation,
--     s.horizontal_datum_id,
--     s.vertical_datum_id,
--     v.name AS vertical_datum,
--     s.huc,
--     s.state_abbrev,
--     st.name as state,
--     COALESCE(code_agg.parameter_codes, '{}') as parameter_codes,
--     json_build_object('action', ns.action, 'flood', ns.flood, 
--     	'moderate_flood', ns.moderate_flood, 'major_flood', ns.major_flood) AS nws_stages,
--     s.create_date,
--     s.update_date
--     FROM usgs_site s
--     JOIN vertical_datum v ON v.id=s.vertical_datum_id
--     JOIN tiger_data.state_all st ON st.stusps=s.state_abbrev
--     LEFT JOIN nws_stages ns ON ns.usgs_site_number=s.site_number
--     LEFT JOIN (
--         SELECT array_agg(code) AS parameter_codes, b.site_id 
--         FROM usgs_parameter a
--         JOIN usgs_site_parameters b ON b.parameter_id = a.id
--         GROUP BY b.site_id
--         ) code_agg ON code_agg.site_id = s.id
-- );

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