CREATE extension IF NOT EXISTS "uuid-ossp";

----------
-- DOMAINS
----------

-- config (application config variables)
CREATE TABLE IF NOT EXISTS config (
    config_name VARCHAR UNIQUE NOT NULL,
    config_value VARCHAR NOT NULL
);

INSERT INTO config (config_name, config_value) VALUES
('write_to_bucket', 'castle-data-develop');

-- Create vertical_datum table
CREATE TABLE IF NOT EXISTS vertical_datum (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR UNIQUE NOT NULL
);

INSERT into vertical_datum (id, name) VALUES
    (0, 'UNKNOWN'),
    (1, 'COE1912'),
    (2, 'NGVD29'),
    (3, 'NAVD88');

CREATE TABLE IF NOT EXISTS location_kind (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL
);

INSERT INTO location_kind (id, name) VALUES
	('83726bc6-82ca-423b-97d2-0309bee76fa7','SITE'),
	('460ea73b-c65e-4fc8-907a-6e6fd2907a99','PROJECT'),
	('c5841fa2-f6cf-4feb-abd6-48798b6cbd48','STREAM'),
	('1e77acaf-fdee-4e7c-b659-101bce76a229','STREAM_LOCATION'),
	('a8ab21ac-ca4c-48d3-8f84-27d920faad14','PUMP'),
	('3598d16a-7ac3-4b94-8f7e-cf73d3ac03c7','EMBANKMENT'),
	('b9494fd2-7504-4412-8216-b38a6d1d0552','OUTLET'),
	('d61950ae-7df0-4fd2-8f03-570ac9fd23bf','GATE'),
	('6bfa45de-1a20-48ba-9ca6-ab267b5e81c7','BASIN'),
	('c7dfa48d-f601-4489-a96b-3d51406e4701','LOCK'),
	('b4c79a9c-0c3f-4934-9270-a611489ed17b','TURBINE'),
	('78a97f08-de50-49e6-b664-5fded7dcb490','STREAM_REACH');

------------
-- LOCATIONS
------------

CREATE TABLE IF NOT EXISTS location (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    office_id UUID,
    name VARCHAR,
    public_name VARCHAR,
    slug VARCHAR UNIQUE NOT NULL,
    geometry geometry,
    kind_id UUID NOT NULL REFERENCES location_kind (id),
    creator UUID,
    create_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    updater UUID,
    update_date TIMESTAMPTZ,
    CONSTRAINT office_unique_name UNIQUE(office_id,name)
);

-- Partial Location Seed data for testing
INSERT INTO location (office_id, name, public_name, slug, geometry, kind_id) VALUES
-- Sample of LRH
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Adamsville','Adamsville','adamsville-1',ST_GeomFromText('POINT(-82.3558 38.8736)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Alderson','Alderson','alderson',ST_GeomFromText('POINT(-80.6413 37.7243)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Alexandria','Alexandria','alexandria-1',ST_GeomFromText('POINT(-82.6106 40.0869)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Allisonia','Allisonia','allisonia',ST_GeomFromText('POINT(-80.7458 36.9375)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr','Alum Creek Lake','alumcr',ST_GeomFromText('POINT(-82.9668 40.18572)',4326),'460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-City Pump Station','City of Columbus Water Supply','alumcr-city-pump-station',ST_GeomFromText('POINT(-82.9652 40.18437)',4326),'a8ab21ac-ca4c-48d3-8f84-27d920faad14'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-DELCO Water Treatment Plant','Delaware County Water Supply','alumcr-delco-water-treatment-plant',ST_GeomFromText('POINT(-82.96655 40.18384)',4326),'a8ab21ac-ca4c-48d3-8f84-27d920faad14'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Dam','Alum Creek Dam','alumcr-dam',ST_GeomFromText('POINT(-82.96389 40.18417)',4326),'3598d16a-7ac3-4b94-8f7e-cf73d3ac03c7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Delco','Delco Water Pump','alumcr-delco',ST_GeomFromText('POINT(-82.9668 40.18572)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Lake','Alum Creek Lake','alumcr-lake',ST_GeomFromText('POINT(-82.9668 40.1858)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Left Spillway','Left Spillway','alumcr-left-spillway',ST_GeomFromText('POINT(-82.96167 40.1841)',4326),'b9494fd2-7504-4412-8216-b38a6d1d0552'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Middle Spillway','Middle Spillway','alumcr-middle-spillway',ST_GeomFromText('POINT(-82.96167 40.1841)',4326),'b9494fd2-7504-4412-8216-b38a6d1d0552'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Outflow','Alum Creek Outflow','alumcr-outflow',ST_GeomFromText('POINT(-82.96149 40.18228)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Right Spillway','Right Spillway','alumcr-right-spillway',ST_GeomFromText('POINT(-82.96167 40.1841)',4326),'b9494fd2-7504-4412-8216-b38a6d1d0552'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr-Water Treatment Plant','Westerville Water Treatment','alumcr-water-treatment-plant',ST_GeomFromText('POINT(-82.96167 40.1841)',4326),'a8ab21ac-ca4c-48d3-8f84-27d920faad14'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','AlumCr_R','Alum Creek','alumcr-r',ST_GeomFromText('POINT(-82.90724 39.8814)',4326),'c5841fa2-f6cf-4feb-abd6-48798b6cbd48'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Ashford','Ashford','ashford',ST_GeomFromText('POINT(-81.71167 38.17972)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Ashland','Ashland','ashland-1',ST_GeomFromText('POINT(-82.6375 38.4814)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Athens','Athens','athens-1',ST_GeomFromText('POINT(-82.0877 39.3288)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood','Atwood Lake','atwood-1',ST_GeomFromText('POINT(-81.285 40.52667)',4326),'460ea73b-c65e-4fc8-907a-6e6fd2907a99'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood Siphon',null,'atwood-siphon',ST_GeomFromText('POINT(0 0)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood-Dam','Atwood Dam','atwood-dam',ST_GeomFromText('POINT(-81.285 40.52667)',4326),'3598d16a-7ac3-4b94-8f7e-cf73d3ac03c7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood-Lake','Atwood Lake','atwood-lake-1',ST_GeomFromText('POINT(-81.285 40.5267)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood-Ouflow',null,'atwood-ouflow',ST_GeomFromText('POINT(-81.285 40.52667)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood-Outflow','Atwood Outflow','atwood-outflow',ST_GeomFromText('POINT(-81.2906 40.5244)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood-Siphon',null,'atwood-siphon-1',ST_GeomFromText('POINT(-81.285 40.52667)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Atwood-Spillway','Spillway','atwood-spillway',ST_GeomFromText('POINT(-81.285 40.52667)',4326),'b9494fd2-7504-4412-8216-b38a6d1d0552'),
    -- Sample of LRN
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ACST1','Sycamore Creek near Ashland City, TN','acst1',ST_GeomFromText('POINT(-87.05111 36.32)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ACST1-SycamoreCr-AshlandCityTN','Sycamore Creek near Ashland City, TN','acst1-sycamorecr-ashlandcitytn',ST_GeomFromText('POINT(-87.05111 36.32)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','AKLN6','Arkport Lake near Arkport NY','akln6',ST_GeomFromText('POINT(-77.71667 42.39583)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ALBK2','Albany KY Rain Gage','albk2',ST_GeomFromText('POINT(-85.12167 36.74528)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ALBK2-AlbanyKY','Albany KY Rain Gage','albk2-albanyky',ST_GeomFromText('POINT(-85.12167 36.74528)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ALPT1','West Fork Obey River near Alpine, TN','alpt1',ST_GeomFromText('POINT(-85.1744 36.39722)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ALPT1-WFkObeyR-AlpineTN','West Fork Obey River near Alpine, TN','alpt1-wfkobeyr-alpinetn',ST_GeomFromText('POINT(-85.1744 36.39722)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','AMOK2','Clarks River - Almo KY','amok2',ST_GeomFromText('POINT(-88.27365 36.69172)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','AMOK2-ClarksR-AlmoKY','Clarks River - Almo KY','amok2-clarksr-almoky',ST_GeomFromText('POINT(-88.27361 36.69167)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ANTT1','Mill Creek near Antioch, TN','antt1',ST_GeomFromText('POINT(-86.68073 36.0816)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ANTT1-MillCr-AntiochTN','Mill Creek near Antioch, TN','antt1-millcr-antiochtn',ST_GeomFromText('POINT(-86.68083 36.08167)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','APLV2','Appalachia VA Rain Gage','aplv2',ST_GeomFromText('POINT(-82.78778 36.90167)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','APLV2-AppalachiaVA','Appalachia VA Rain Gage','aplv2-appalachiava',ST_GeomFromText('POINT(-82.78778 36.90167)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ARTT1','Powell River - Arthur TN','artt1',ST_GeomFromText('POINT(-83.63028 36.54194)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ARTT1-PowellR-arthurTN','Powell River - Arthur TN','artt1-powellr-arthurtn',ST_GeomFromText('POINT(-83.63028 36.54194)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ASHT1','Cheatham Dam Tailwater','asht1',ST_GeomFromText('POINT(-87.22542 36.32042)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ASHT1-ASHT1-CHEATHAM',null,'asht1-asht1-cheatham',ST_GeomFromText('POINT(-87.22542 36.32042)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','ASHT1-CHEATHAM','Cheatham Dam Tailwater','asht1-cheatham',ST_GeomFromText('POINT(-87.22542 36.32042)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('91cf44dc-6384-4622-bd9f-0f36e4343413','DaleHollow','DALE HOLLOW LAKE NEAR CELINA','dalehollow',ST_GeomFromText('POINT(-85.45639 36.54139)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
	('91cf44dc-6384-4622-bd9f-0f36e4343413','DaleHollow-Lake','DALE HOLLOW LAKE NEAR CELINA','dalehollow-lake',ST_GeomFromText('POINT(-85.45639 36.54139)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
	('91cf44dc-6384-4622-bd9f-0f36e4343413','DaleHollow-Tailwater','DALE HOLLOW TAILWATER','dalehollow-tailwater',ST_GeomFromText('POINT(-85.46111 36.54722)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    -- Sample of MVP
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ABBW3S','Abbotsford Snow','abbw3s',ST_GeomFromText('POINT(-90.25944 44.945)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ABRN8','Wild Rice River at Abercrombie','abrn8',ST_GeomFromText('POINT(-96.7825 46.47028)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ABRN8S','Abercrombie 3NW Snow','abrn8s',ST_GeomFromText('POINT(-96.79 46.48)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ADAM5S','Ada Snow','adam5s',ST_GeomFromText('POINT(-96.54306 47.29389)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ADMN8S','Adams Snow','admn8s',ST_GeomFromText('POINT(-98.0775 48.44528)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ADSM5','Wild Rice River near Ada','adsm5',ST_GeomFromText('POINT(-97.05 47.2639)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AGYM5','Middle River at Argyle','agym5',ST_GeomFromText('POINT(-96.81611 48.34028)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AHLM5S','Ash Lake Snow','ahlm5s',ST_GeomFromText('POINT(-92.92871 48.23399)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AKRN8','Tongue River at Akra','akrn8',ST_GeomFromText('POINT(-97.74639 48.77833)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ALVM5','Snake River below Alvarado','alvm5',ST_GeomFromText('POINT(-96.7764 48.1958)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AMAW3','Mississippi River near Alma','amaw3',ST_GeomFromText('POINT(-91.9 44.2833)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AMEN8','Rush River nr Amenia','amen8',ST_GeomFromText('POINT(-97.21361 47.01528)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AMEN8S','Amenia Snow','amen8s',ST_GeomFromText('POINT(-97.22063 47.00593)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ANEN8S','ANETA Snow','anen8s',ST_GeomFromText('POINT(-97.98806 47.69389)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ANKM5','Mississippi River near Anoka','ankm5',ST_GeomFromText('POINT(-93.2967 45.1267)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ANNM5S','Annandale Snow','annm5s',ST_GeomFromText('POINT(-94.10306 45.25389)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ANTW3S','Anitigo Snow','antw3s',ST_GeomFromText('POINT(-89.14972 45.18472)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','APPM5','Pomme de Terre River Appleton','appm5',ST_GeomFromText('POINT(-96.0194 45.2039)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','APPM5S','Appleton Snow','appm5s',ST_GeomFromText('POINT(-96.0201 45.2037)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ATKM5','Miss River near Aitkin, MN','atkm5',ST_GeomFromText('POINT(-93.7072 46.5406)',4326),'1e77acaf-fdee-4e7c-b659-101bce76a229'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','ATKM5S','Aitkin Snow','atkm5s',ST_GeomFromText('POINT(-93.73383 46.52283)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','Adams','Adams Soil Moisture/Temp Gage','adams',ST_GeomFromText('POINT(-98.07587 48.49981)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','AitkinDiversion','Mississippi River Aitkin, MN','aitkindiversion',ST_GeomFromText('POINT(93.41 46.35)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7');

-- Location regression test
INSERT INTO location (id, office_id, name, public_name, slug, geometry, kind_id) VALUES
	('45cd0d9f-6751-434f-afe4-9da0690793be','91cf44dc-6384-4622-bd9f-0f36e4343413','RegressionTestLocation01','Regression Test Location 01','regression-test-location-01',ST_GeomFromText('POINT(-82.35583 38.87361)',4326),'83726bc6-82ca-423b-97d2-0309bee76fa7');


------------
-- LEVELS
------------

-- level_kind definition
CREATE TABLE IF NOT EXISTS level_kind (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR UNIQUE NOT NULL
);

-- level definition
CREATE TABLE IF NOT EXISTS level (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    location_id UUID,
    level_kind_id UUID NOT NULL REFERENCES level_kind(id),
    CONSTRAINT unique_location_level_kind UNIQUE(location_id, level_kind_id)
);

-- level_value definition
CREATE TABLE IF NOT EXISTS level_value (
    level_id UUID NOT NULL REFERENCES level(id),
    julian_date DOUBLE PRECISION NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    CONSTRAINT unique_level_julian UNIQUE(level_id, julian_date)
);

------------
-- OFFICE
------------

CREATE TABLE IF NOT EXISTS office (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL,
    symbol VARCHAR UNIQUE NOT NULL,
    parent_id UUID references office(id)
);

INSERT INTO office (id, name, symbol, parent_id) VALUES
    ('91cf44dc-6384-4622-bd9f-0f36e4343413','Great Lakes and Ohio River Division','LRD',Null),
    ('17fa25b8-44a0-4e6d-9679-bdf6b0ee6b1a','Buffalo District','LRB','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('d8f8934d-e414-499d-bd51-bc93bbde6345','Chicago District','LRC','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('a8192ad1-206c-4da6-b19e-b7ba7a67aa1f','Detroit District','LRE','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Huntington District','LRH','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('433a554d-7b27-4046-89eb-906788eb4046','Louisville District','LRL','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','Nashville District','LRN','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('61291eaf-d62f-4846-ad95-87cc86b56851','Pittsburgh District','LRP','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('485d800d-a30d-4fcb-af43-0bea2ce11adb','Mississippi Valley Division','MVD',Null),
    ('1245e3c0-fc72-4621-86b2-24ff7de21f88','Memphis District','MVM','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('f81f5659-ce57-4c87-9c7a-0d685a983bfd','New Orleans District','MVN','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('81557734-7046-4c55-90ac-066dd882166a','Rock Island District','MVR','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('565be474-0c68-44a6-8e66-b833a39685bd','St. Louis District','MVS','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','St. Paul District','MVP','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('b9cca282-eb91-4ea1-b075-d067b4420184','Vicksburg District','MVK','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('973ad07b-7df3-4a95-9e43-7bc25930f7a8','North Atlantic Division','NAD',Null),
    ('7ed4821f-9e37-4c56-8baf-05c1b5bcc84c','Baltimore District','NAB','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('bf99dc79-51d3-4abe-aba4-7e1781315766','Europe District','NAU','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('30cb06ee-bd94-4c49-a945-e92735e7bdc1','New England District','NAE','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('4ca9e255-8a88-44d3-8091-bb61931e600c','New York District','NAN','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('a47f1ef4-1017-43c1-bf36-67f57376d163','Norfolk District','NAO','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('1989e3fc-f12a-42da-a263-c3ae978e2c09','Philadelphia District','NAP','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('ad713a67-37d6-444e-b6b6-e02c0858451f','Northwestern Division','NWD',Null),
    ('5b35ea7c-8d1b-481a-956b-b32939093db4','Kansas City District','NWK','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('665ffb00-ccba-488c-83c5-2083543cacd7','Omaha District','NWO','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('8b0a732d-d511-4332-b2e7-dd6943a597e9','Portland District','NWP','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('007cbff5-6946-4b9b-a3f7-0bef4406f122','Seattle District','NWS','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('30266178-d32a-4e07-aea1-1f7182ed245e','Walla Walla District','NWW','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('ab99f33f-836e-4788-a931-33e0376d1406','Pacific Ocean Division','POD',Null),
    ('be0614bf-1461-4993-9ce7-8d1d17606be9','Alaska District','POA','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('64cd2766-2586-4709-a4b9-f8a6029233ea','Far East District','POF','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('8b86f8cb-0594-4d69-a66c-e4e295c2b5af','Honolulu District','POH','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('f7300efc-ff48-44fd-b43f-b5373a42ba3e','Japan Engineer District','POJ','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('e046baab-c0b6-4dcf-8cc7-cbab155dddc0','South Atlantic Division','SAD',Null),
    ('d4501358-1c48-45cb-86f3-f1a31e9bd93f','Charleston District','SAC','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('b4f45596-70e5-4a12-a894-a64300648244','Jacksonville District','SAJ','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('9ffc189c-ad40-4fbf-bc06-2098c6cb920e','Mobile District','SAM','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('0154184e-2509-4485-b449-8eff4ab52eef','Savannah District','SAS','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('ba1f7846-43d9-4a21-9876-27c59510d9c0','Wilmington District','SAW','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('f92cb397-2c8c-44f1-856d-a00ef9467224','South Pacific Division','SPD',Null),
    ('11b5fe49-fe36-4a06-a0da-d55b1b62b1fb','Albuquerque District','SPA','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('b8cec5bc-f975-4bed-993d-8f913ca51673','Los Angeles District','SPL','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('ff52a84b-356a-4173-a8df-89a1b408d354','Sacramento District','SPK','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('cf9b1f4d-1cd3-4a00-b73d-b6f8fe75915e','San Francisco District','SPN','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('2176fa5b-7d6f-4f73-8dc3-18e2323aefbb','Southwestern Division','SWD',Null),
    ('f3f0d7ff-19b6-4167-a3f1-5c04f5a0fe4d','Fort Worth District','SWF','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    ('72ee5695-cdaa-4182-b0c1-4d27f1a3f570','Galveston District','SWG','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    ('131daa5c-49c2-4488-be6b-bd638a83a03f','Little Rock District','SWL','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    ('fe29f6e2-e200-44a4-9545-a4680ab9366e','Tulsa District','SWT','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb');

------------
-- USGS_SITE
------------

-- Create usgs_site table
CREATE TABLE IF NOT EXISTS usgs_site (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    site_number VARCHAR UNIQUE NOT NULL,
    name VARCHAR,
    geometry geometry,
    elevation REAL,
    horizontal_datum_id INTEGER NOT NULL REFERENCES public.spatial_ref_sys(srid),
    vertical_datum_id INTEGER NOT NULL REFERENCES vertical_datum(id), 
    huc VARCHAR,
    state_abbrev VARCHAR(2) NOT NULL,
    create_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_date TIMESTAMPTZ
);

-- Create usgs_parameter table
CREATE TABLE IF NOT EXISTS usgs_parameter (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    code VARCHAR UNIQUE NOT NULL,
    description VARCHAR NOT NULL
);

-- usgs_parameter seed data
INSERT INTO usgs_parameter (id, code, description) VALUES 
('a9f78261-e6a6-4ad2-827e-bd7a4ac0dc28', '00065', 'Gage height, feet'),
('ba29fc34-6315-4424-838f-9b1863715fad', '00060', 'Discharge, cubic feet per second'),
('06cca640-f52b-4c0c-8f64-a985836fda5a', '00061', 'Discharge, cubic feet per second, instantaneous'),
('f739b4af-1c96-437c-a788-901f59d177fb', '62614', 'Lake or reservoir water surface elevation above NGVD 1929, feet'),
('60bb26cb-a65d-40d2-ad54-b00d6802de7b', '62615', 'Lake or reservoir water surface elevation above NAVD 1988, feet'),
('738eb4df-b34b-45cc-a5aa-f2136384244f', '00045', 'Precipitation, total, inches'),
('0fa9993d-6674-4ba3-ac8a-f02830beea1e', '00010', 'Temperature, water, degrees Celsius'),
('12ff9f0b-159b-43cb-8126-5253f7948380', '00011', 'Temperature, water, degrees Fahrenheit');

-- Create usgs_site_parameters table
CREATE TABLE IF NOT EXISTS usgs_site_parameters (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    site_id UUID NOT NULL REFERENCES usgs_site(id),
    parameter_id UUID NOT NULL REFERENCES usgs_parameter(id),
    CONSTRAINT site_unique_param UNIQUE(site_id, parameter_id)
);

-- usgs_measurements
CREATE TABLE IF NOT EXISTS usgs_measurements (
    time TIMESTAMPTZ NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    usgs_site_parameters_id UUID NOT NULL REFERENCES usgs_site_parameters (id) ON DELETE CASCADE,
    CONSTRAINT site_parameters_unique_time UNIQUE(usgs_site_parameters_id, time),
    PRIMARY KEY (usgs_site_parameters_id, time)
);

------------
-- WATERSHED
------------

-- watershed
CREATE TABLE IF NOT EXISTS watershed (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    geometry geometry NOT NULL DEFAULT ST_GeomFromText('POLYGON ((0 0, 0 0, 0 0, 0 0, 0 0))',5070),
    office_id UUID NOT NULL REFERENCES office(id),
	deleted boolean NOT NULL DEFAULT false
);

-- extent to polygon reference order - simple 4 point extents
-- xmin,ymax (top left), xmax ymax (top right), xmax ymin (bottom right), xmin ymin (bottom left), xmin ymax (top left again)
INSERT INTO watershed (id,slug,"name",geometry,office_id) VALUES	 
    ('0f065e6a-3380-4ac3-b576-89fae7774b9f','little-sandy-river','Little Sandy River',ST_GeomFromText('POLYGON ((1096000 1812000, 1158000 1812000, 1158000 1732000, 1096000 1732000, 1096000 1812000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),    
    ('1a629fac-82c9-4b3e-b7fc-6a891d944140','ohio-river','Ohio River',ST_GeomFromText('POLYGON ((1006000 1914000, 1206000 1914000, 1206000 1754000, 1006000 1754000, 1006000 1914000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	
    ('3e322a11-b76b-4710-8f9a-b7884cd8ae77','big-sandy-river','Big Sandy River',ST_GeomFromText('POLYGON ((1114000 1796000, 1288000 1796000, 1288000 1624000, 1114000 1624000, 1114000 1796000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('4d3083d1-101c-4b76-9311-1154917ffbf1','twelvepole-river','Twelvepole River',ST_GeomFromText('POLYGON ((1152000 1796000, 1212000 1796000, 1212000 1728000, 1152000 1728000, 1152000 1796000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('5024720e-02f6-4577-a09c-ff1ff5a28223','hocking-river','Hocking River',ST_GeomFromText('POLYGON ((1112000 1960000, 1220000 1960000, 1220000 1878000, 1112000 1878000, 1112000 1960000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('50372dbc-f254-4584-8345-1c3613d2a102','guyandotte-river','Guyandotte River',ST_GeomFromText('POLYGON ((1166000 1814000, 1298000 1814000, 1298000 1692000, 1166000 1692000, 1166000 1814000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('5758d0dc-c8bf-4e37-a5e7-44ff3f4b8677','scioto-river','Scioto River',ST_GeomFromText('POLYGON ((1004000 2056000, 1154000 2056000, 1154000 1810000, 1004000 1810000, 1004000 2056000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('65a93467-c9b4-4166-acb6-58e8ec06ed3b','kanawha-river','Kanawha River',ST_GeomFromText('POLYGON ((1182000 1870000, 1410000 1870000, 1410000 1544000, 1182000 1544000, 1182000 1870000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('7c6dd902-fbc5-43e4-9bbf-351963f5723d','muskingum-river','Muskingum River',ST_GeomFromText('POLYGON ((1098000 2110000, 1268000 2110000, 1268000 1904000, 1098000 1904000, 1098000 2110000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('cf193b4e-61c3-4e4d-9503-2935a82aed96','little-kanawha-river','Little Kanawha River',ST_GeomFromText('POLYGON ((1164000 1970000, 1354000 1970000, 1354000 1824000, 1164000 1824000, 1164000 1970000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
	('c54eab5b-1020-476b-a5f8-56d77802d9bf','tennessee-river','Tennessee River',ST_GeomFromText('POLYGON ((640000 1678000, 1300000 1678000, 1300000 1268000, 640000 1268000, 640000 1678000))',5070),'552e59f7-c0cc-4689-8a4d-e791c028430a'),	 
    ('c785f4de-ab17-444b-b6e6-6f1ad16676e8','cumberland-basin-river','Cumberland Basin River',ST_GeomFromText('POLYGON ((662000 1678000, 1172000 1678000, 1172000 1408000, 662000 1408000, 662000 1678000))',5070),'552e59f7-c0cc-4689-8a4d-e791c028430a'),	 																																	
	('feda585b-1ba0-4b19-92ed-7195154b8052','tennessee-cumberland-river', 'Tennessee & Cumberland River', ST_GeomFromText('POLYGON ((642000 1682000, 1300000 1682000, 1300000 1258000, 642000 1258000, 642000 1682000))',5070), '552e59f7-c0cc-4689-8a4d-e791c028430a'),
    ('03206ff6-fe91-426c-a5e9-4c651b06f9c6','eau-galla-river','Eau Galla River',ST_GeomFromText('POLYGON ((284000 2460000, 326000 2460000, 326000 2404000, 284000 2404000, 284000 2460000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
    ('048ce853-6642-4ac4-9fb2-81c01f67a85b','mississippi-river-headwaters','Mississippi River Headwaters',ST_GeomFromText('POLYGON ((24000 2760000, 254000 2760000, 254000 2402000, 24000 2402000, 24000 2760000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
    ('ad30f178-afc3-43b9-ba92-7bd139581217','red-river-north','Red River North',ST_GeomFromText('POLYGON ((-356000 2950000, 150000 2950000, 150000 2494000, -356000 2494000, -356000 2950000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
    ('c8bf6c6d-7f19-406a-a438-f2f876ce4815','souris-river','Souris River',ST_GeomFromText('POLYGON ((-708000 3100000, -178000 3100000, -178000 2736000, -708000 2736000, -708000 3100000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
    ('ced6ec9e-43b5-496e-a2b7-894af92c9b63','mississippi-river-navigation','Mississippi River Navigation',ST_GeomFromText('POLYGON ((48000 2646000, 564000 2646000, 564000 2204000, 48000 2204000, 48000 2646000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
    ('f4219691-e498-46a3-ab0f-f2957bd09a10','minnesota-river','Minnesota River',ST_GeomFromText('POLYGON ((-112000 2602000, 234000 2602000, 234000 2244000, -112000 2244000, -112000 2602000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321')
    ('c572ed70-d401-4a97-aea6-cb3fe2b77e41','savannah-river-basin','Savannah River Basin',ST_GeomFromText('POLYGON ((1110000 1432000, 1432000 1432000, 1432000 1094000, 1110000 1094000, 1110000 1432000))',5070),'0154184e-2509-4485-b449-8eff4ab52eef');

-- watershed_usgs_sites
CREATE TABLE IF NOT EXISTS watershed_usgs_sites (
    watershed_id UUID NOT NULL REFERENCES watershed(id),
    usgs_site_parameter_id UUID NOT NULL REFERENCES usgs_site_parameters(id),
    CONSTRAINT watershed_unique_site_param UNIQUE(watershed_id, usgs_site_parameter_id)
);
-- Add comment to describe table
COMMENT ON TABLE watershed_usgs_sites IS 'This is a bridge table.  Each entry represent a watershed/site/parameter that has been requested for USGS data acquisition.';

-------------
-- NWS STAGES
-------------

-- Create usgs_site table
CREATE TABLE IF NOT EXISTS nws_stages (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    nwsid VARCHAR UNIQUE NOT NULL,
    usgs_site_number VARCHAR UNIQUE NOT NULL REFERENCES usgs_site(site_number),
    name VARCHAR,
    action DOUBLE PRECISION,
    flood DOUBLE PRECISION,
    moderate_flood DOUBLE PRECISION,
    major_flood DOUBLE PRECISION,
    create_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_date TIMESTAMPTZ
);

-------------------
-- SHAPEFILE UPLOAD
-------------------

-- upload_status definition
CREATE TABLE IF NOT EXISTS upload_status (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL
);

INSERT INTO upload_status (id, name) VALUES
    ('b5d777fc-c46b-4a10-a488-1415e3d7849d', 'INITIATED'),
    ('969e00ad-2be1-4cf5-9f80-5c198e1e8450', 'SUCCESS'),
    ('020c8cda-913b-4c2d-8580-2834381bf885', 'FAIL');


-- watershed_shapefile_uploads definition
CREATE TABLE IF NOT EXISTS watershed_shapefile_uploads (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	watershed_id UUID NOT NULL REFERENCES watershed(id),
	file VARCHAR NOT NULL,
	date_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	file_size INTEGER NOT NULL,
    processing_info VARCHAR,
    user_id UUID,
    upload_status_id UUID NOT NULL DEFAULT 'b5d777fc-c46b-4a10-a488-1415e3d7849d' REFERENCES upload_status(id)
);