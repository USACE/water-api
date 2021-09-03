-- usgs_measurements
CREATE TABLE IF NOT EXISTS usgs_measurements (
    time TIMESTAMPTZ NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    usgs_site_parameters_id UUID NOT NULL REFERENCES usgs_site_parameters (id) ON DELETE CASCADE,
    CONSTRAINT site_parameters_unique_time UNIQUE(usgs_site_parameters_id, time),
    PRIMARY KEY (usgs_site_parameters_id, time)
);

-- Grant read
GRANT SELECT ON usgs_measurements TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON usgs_measurements TO water_writer;


-- Sample Stage values for GUYANDOTTE RIVER AT LOGAN, WV (03203600)
INSERT INTO usgs_measurements(time, value, usgs_site_parameters_id) VALUES
('2021-08-27 00:00:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 00:15:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 00:30:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 00:45:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 01:00:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 01:15:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 01:30:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 01:45:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 02:00:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 02:15:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 02:30:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 02:45:00', 4.88, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 03:00:00', 4.85, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 03:15:00', 4.84, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 03:30:00', 4.84, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 03:45:00', 4.84, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 04:00:00', 4.84, '2a8c983a-2210-490b-a18d-55533a048f4a'),
('2021-08-27 04:15:00', 4.84, '2a8c983a-2210-490b-a18d-55533a048f4a');