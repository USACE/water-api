-- Create usgs_parameter table
CREATE TABLE IF NOT EXISTS usgs_parameter (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,
    description VARCHAR NOT NULL
);

-- Create usgs_site_parameters table
CREATE TABLE IF NOT EXISTS usgs_site_parameters (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    usgs_site_id UUID NOT NULL REFERENCES usgs_site(id),
    usgs_parameter_id UUID NOT NULL REFERENCES usgs_parameter(id),
    CONSTRAINT site_unique_param UNIQUE(usgs_site_id, usgs_parameter_id)
);

-- Create v_usgs_site_parameters_enabled
CREATE OR REPLACE VIEW v_usgs_site_parameters_enabled AS (
    SELECT
    distinct usp.usgs_site_id,
    us.usgs_id,
    code_agg.usgs_parameter_code,
    param_id_agg.usgs_parameter_id
    FROM a2w_cwms.usgs_site_parameters usp
    JOIN a2w_cwms.usgs_site us ON us.id = usp.usgs_site_id
    INNER JOIN (
        SELECT array_agg(usgs_parameter_id) AS usgs_parameter_id, a.usgs_site_id 
        FROM a2w_cwms.usgs_site_parameters a
        JOIN a2w_cwms.usgs_parameter b ON a.usgs_parameter_id = b.id
        GROUP BY a.usgs_site_id 
        ) param_id_agg ON param_id_agg.usgs_site_id = us.id	
    INNER JOIN (
        SELECT array_agg(code) AS usgs_parameter_code, b.usgs_site_id 
        FROM a2w_cwms.usgs_parameter a
        JOIN a2w_cwms.usgs_site_parameters b ON b.usgs_parameter_id = a.id
        GROUP BY b.usgs_site_id 
        ) code_agg ON code_agg.usgs_site_id = us.id	
    ORDER BY us.usgs_id
);

-- Grant read
GRANT SELECT ON
    usgs_parameter,
    usgs_site_parameters,
    v_usgs_site_parameters_enabled
TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON
    usgs_parameter,
    usgs_site_parameters
TO water_writer;

-- usgs_parameter seed data
INSERT INTO usgs_parameter (id, code, description) VALUES 
('a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28', '00065', 'Gage height, feet'),
('ba29fc34-6315-4424-838f-9b1863715fad', '00060', 'Discharge, cubic feet per second'),
('f739b4af-1c96-437c-a788-901f59d177fb', '62614', 'Lake or reservoir water surface elevation above NGVD 1929, feet'),
('60bb26cb-a65d-40d2-ad54-b00d6802de7b', '62615', 'Lake or reservoir water surface elevation above NAVD 1988, feet'),
('738eb4df-b34b-45cc-a5aa-f2136384244f', '00045', 'Precipitation, total, inches'),
('0fa9993d-6674-4ba3-ac8a-f02830beea1e', '00010', 'Temperature, water, degrees Celsius'),
('12ff9f0b-159b-43cb-8126-5253f7948380', '00011', 'Temperature, water, degrees Fahrenheit');

-- usgs_site_parameters seed data for testing
INSERT INTO usgs_site_parameters (id, usgs_site_id, usgs_parameter_id) VALUES
-- GUYANDOTTE RIVER AT LOGAN, WV - Stage and Flow
('2a8c983a-2210-490b-a18d-55533a048f4a', (select id from usgs_site where name='GUYANDOTTE RIVER AT LOGAN, WV'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('b5ad3c36-4238-4fbb-8b0d-a5d544479eac', (select id from usgs_site where name='GUYANDOTTE RIVER AT LOGAN, WV'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
-- GUYANDOTTE RIVER AT BRANCHLAND, WV - Stage and Precip
('1fdd9651-84ab-4d17-9e6f-37a5c2596521', (select id from usgs_site where name='GUYANDOTTE RIVER AT BRANCHLAND, WV'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('ef9538de-8e44-4827-b552-0498ef1c18ff', (select id from usgs_site where name='GUYANDOTTE RIVER AT BRANCHLAND, WV'), '738eb4df-b34b-45cc-a5aa-f2136384244f');