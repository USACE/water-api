-- watershed_usgs_sites
CREATE TABLE IF NOT EXISTS watershed_usgs_sites (
    watershed_uid UUID REFERENCES watershed(uid),
    usgs_site_uid UUID REFERENCES usgs_site(uid),
    usgs_parameter_uid UUID REFERENCES usgs_parameter(uid)
);

-- Grant read
GRANT SELECT ON
    watershed_usgs_sites
TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON
    watershed_usgs_sites
TO water_writer;

INSERT into watershed_usgs_sites (watershed_uid, usgs_site_uid, usgs_parameter_uid) VALUES
-- LRH - EAST FORK TWELVEPOLE CREEK NEAR DUNLOW, WV (stage, flow)
('4d3083d1-101c-4b76-9311-1154917ffbf1', (select uid from usgs_site where site_number = '03206600'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('4d3083d1-101c-4b76-9311-1154917ffbf1', (select uid from usgs_site where site_number = '03206600'), 'ba29fc34-6315-4424-838f-9b1863715fad'),
-- LRH - TWELVEPOLE CREEK BELOW WAYNE, WV (stage, precip)
('4d3083d1-101c-4b76-9311-1154917ffbf1', (select uid from usgs_site where site_number = '03207020'), 'a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28'),
('4d3083d1-101c-4b76-9311-1154917ffbf1', (select uid from usgs_site where site_number = '03207020'), '738eb4df-b34b-45cc-a5aa-f2136384244f');