-- Create usgs_parameter table
CREATE TABLE IF NOT EXISTS usgs_parameter (
    uid UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,
    description VARCHAR NOT NULL
);

-- Grant read
GRANT SELECT ON usgs_parameter TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON usgs_parameter TO water_writer;

-- usgs_parameter seed data
INSERT INTO usgs_parameter (uid, code, description) VALUES 
('a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28', '00065', 'Gage height, feet'),
('ba29fc34-6315-4424-838f-9b1863715fad', '00060', 'Discharge, cubic feet per second'),
('f739b4af-1c96-437c-a788-901f59d177fb', '62614', 'Lake or reservoir water surface elevation above NGVD 1929, feet'),
('60bb26cb-a65d-40d2-ad54-b00d6802de7b', '62615', 'Lake or reservoir water surface elevation above NAVD 1988, feet'),
('738eb4df-b34b-45cc-a5aa-f2136384244f', '00045', 'Precipitation, total, inches'),
('0fa9993d-6674-4ba3-ac8a-f02830beea1e', '00010', 'Temperature, water, degrees Celsius'),
('12ff9f0b-159b-43cb-8126-5253f7948380', '00011', 'Temperature, water, degrees Fahrenheit');
