-- Create usgs_site_parameters table
CREATE TABLE IF NOT EXISTS usgs_site_parameters (
    uid UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    site_uid UUID NOT NULL REFERENCES usgs_site(uid),
    parameter_uid UUID NOT NULL REFERENCES usgs_parameter(uid),
    CONSTRAINT site_unique_param UNIQUE(site_uid, parameter_uid)
);

-- Create v_usgs_site
CREATE OR REPLACE VIEW v_usgs_site AS (
    SELECT 
    s.uid,
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
    s.create_date,
    s.update_date
    FROM usgs_site s
    JOIN vertical_datum v ON v.id=s.vertical_datum_id
    JOIN tiger_data.state_all st ON st.stusps=s.state_abbrev
    LEFT JOIN (
        SELECT array_agg(code) AS parameter_codes, b.site_uid 
        FROM usgs_parameter a
        JOIN usgs_site_parameters b ON b.parameter_uid = a.uid
        GROUP BY b.site_uid
        ) code_agg ON code_agg.site_uid = s.uid
);

-- Grant read
GRANT SELECT ON usgs_site_parameters, v_usgs_site TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON usgs_site_parameters TO water_writer;

-- usgs_site_parameters seed data for testing
INSERT INTO usgs_site_parameters (id, usgs_site_id, usgs_parameter_id) VALUES
-- GUYANDOTTE RIVER AT LOGAN, WV - Stage and Flow
('2a8c983a-2210-490b-a18d-55533a048f4a', (select id from usgs_site where name='GUYANDOTTE RIVER AT LOGAN, WV'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('b5ad3c36-4238-4fbb-8b0d-a5d544479eac', (select id from usgs_site where name='GUYANDOTTE RIVER AT LOGAN, WV'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
-- GUYANDOTTE RIVER AT BRANCHLAND, WV - Stage and Precip
('1fdd9651-84ab-4d17-9e6f-37a5c2596521', (select id from usgs_site where name='GUYANDOTTE RIVER AT BRANCHLAND, WV'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('ef9538de-8e44-4827-b552-0498ef1c18ff', (select id from usgs_site where name='GUYANDOTTE RIVER AT BRANCHLAND, WV'), '738eb4df-b34b-45cc-a5aa-f2136384244f')
ON CONFLICT DO NOTHING;