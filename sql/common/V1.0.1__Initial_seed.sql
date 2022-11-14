
INSERT INTO cwms_location_kind (name) VALUES
	('SITE'),
	('PROJECT'),
	('STREAM'),
	('STREAM_LOCATION'),
	('PUMP'),
	('EMBANKMENT'),
	('OUTLET'),
	('GATE'),
	('BASIN'),
	('LOCK'),
	('TURBINE'),
	('STREAM_REACH');

INSERT INTO usgs_site_type (abbreviation, name) VALUES
	('LK', 'Lake, Reservoir, Impoundment'),
    ('ST', 'Stream'),
    ('ST-TS', 'Tidal Stream');

INSERT INTO provider (id, name, slug, parent_id) VALUES
    ('91cf44dc-6384-4622-bd9f-0f36e4343413','Great Lakes and Ohio River Division','lrd',Null),
    ('17fa25b8-44a0-4e6d-9679-bdf6b0ee6b1a','Buffalo District','lrb','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('d8f8934d-e414-499d-bd51-bc93bbde6345','Chicago District','lrc','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('a8192ad1-206c-4da6-b19e-b7ba7a67aa1f','Detroit District','lre','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('2f160ba7-fd5f-4716-8ced-4a29f75065a6','Huntington District','lrh','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('433a554d-7b27-4046-89eb-906788eb4046','Louisville District','lrl','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('552e59f7-c0cc-4689-8a4d-e791c028430a','Nashville District','lrn','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('61291eaf-d62f-4846-ad95-87cc86b56851','Pittsburgh District','lrp','91cf44dc-6384-4622-bd9f-0f36e4343413'),
    ('485d800d-a30d-4fcb-af43-0bea2ce11adb','Mississippi Valley Division','mvd',Null),
    ('1245e3c0-fc72-4621-86b2-24ff7de21f88','Memphis District','mvm','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('f81f5659-ce57-4c87-9c7a-0d685a983bfd','New Orleans District','mvn','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('81557734-7046-4c55-90ac-066dd882166a','Rock Island District','mvr','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('565be474-0c68-44a6-8e66-b833a39685bd','St. Louis District','mvs','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('2cf60156-f22a-418a-bc9f-a28960ad0321','St. Paul District','mvp','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('b9cca282-eb91-4ea1-b075-d067b4420184','Vicksburg District','mvk','485d800d-a30d-4fcb-af43-0bea2ce11adb'),
    ('973ad07b-7df3-4a95-9e43-7bc25930f7a8','North Atlantic Division','nad',Null),
    ('7ed4821f-9e37-4c56-8baf-05c1b5bcc84c','Baltimore District','nab','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('bf99dc79-51d3-4abe-aba4-7e1781315766','Europe District','nau','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('30cb06ee-bd94-4c49-a945-e92735e7bdc1','New England District','nae','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('4ca9e255-8a88-44d3-8091-bb61931e600c','New York District','nan','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('a47f1ef4-1017-43c1-bf36-67f57376d163','Norfolk District','nao','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('1989e3fc-f12a-42da-a263-c3ae978e2c09','Philadelphia District','nap','973ad07b-7df3-4a95-9e43-7bc25930f7a8'),
    ('ad713a67-37d6-444e-b6b6-e02c0858451f','Northwestern Division','nwd',Null),
    ('5b35ea7c-8d1b-481a-956b-b32939093db4','Kansas City District','nwk','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('665ffb00-ccba-488c-83c5-2083543cacd7','Omaha District','nwo','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('8b0a732d-d511-4332-b2e7-dd6943a597e9','Portland District','nwp','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('007cbff5-6946-4b9b-a3f7-0bef4406f122','Seattle District','nws','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('30266178-d32a-4e07-aea1-1f7182ed245e','Walla Walla District','nww','ad713a67-37d6-444e-b6b6-e02c0858451f'),
    ('ab99f33f-836e-4788-a931-33e0376d1406','Pacific Ocean Division','pod',Null),
    ('be0614bf-1461-4993-9ce7-8d1d17606be9','Alaska District','poa','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('64cd2766-2586-4709-a4b9-f8a6029233ea','Far East District','pof','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('8b86f8cb-0594-4d69-a66c-e4e295c2b5af','Honolulu District','poh','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('f7300efc-ff48-44fd-b43f-b5373a42ba3e','Japan Engineer District','poj','ab99f33f-836e-4788-a931-33e0376d1406'),
    ('e046baab-c0b6-4dcf-8cc7-cbab155dddc0','South Atlantic Division','sad',Null),
    ('d4501358-1c48-45cb-86f3-f1a31e9bd93f','Charleston District','sac','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('b4f45596-70e5-4a12-a894-a64300648244','Jacksonville District','saj','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('9ffc189c-ad40-4fbf-bc06-2098c6cb920e','Mobile District','sam','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('0154184e-2509-4485-b449-8eff4ab52eef','Savannah District','sas','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('ba1f7846-43d9-4a21-9876-27c59510d9c0','Wilmington District','saw','e046baab-c0b6-4dcf-8cc7-cbab155dddc0'),
    ('f92cb397-2c8c-44f1-856d-a00ef9467224','South Pacific Division','spd',Null),
    ('11b5fe49-fe36-4a06-a0da-d55b1b62b1fb','Albuquerque District','spa','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('b8cec5bc-f975-4bed-993d-8f913ca51673','Los Angeles District','spl','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('ff52a84b-356a-4173-a8df-89a1b408d354','Sacramento District','spk','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('cf9b1f4d-1cd3-4a00-b73d-b6f8fe75915e','San Francisco District','spn','f92cb397-2c8c-44f1-856d-a00ef9467224'),
    ('2176fa5b-7d6f-4f73-8dc3-18e2323aefbb','Southwestern Division','swd',Null),
    ('f3f0d7ff-19b6-4167-a3f1-5c04f5a0fe4d','Fort Worth District','swf','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    ('72ee5695-cdaa-4182-b0c1-4d27f1a3f570','Galveston District','swg','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    ('131daa5c-49c2-4488-be6b-bd638a83a03f','Little Rock District','swl','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    ('fe29f6e2-e200-44a4-9545-a4680ab9366e','Tulsa District','swt','2176fa5b-7d6f-4f73-8dc3-18e2323aefbb'),
    -- USGS
    ('154791a9-1be9-4b11-a964-3bbd1e08fa89','U.S. Geological Survey','usgs', NULL),
    -- NWS
    ('c3164251-a68a-459f-b0ca-2cdbe982db5a','National Weather Service','noaa-nws', NULL);



INSERT into datatype(id, slug, name, uri) VALUES
    ('4b38762b-d2e0-4842-be78-e9df9c7b8a80', 'cwms-watershed', 'CWMS Watershed', NULL),
    ('a138e363-30ea-4e0d-8d8f-cce03cb8e1d0', 'cwms-timeseries', 'CWMS Timeseries', 'https://cwms-data.usace.army.mil/cwms-data/timeseries'),
    ('97920d27-ee54-4d35-aec4-c01ec31221a2', 'cwms-level', 'CWMS Level', 'https://cwms-data.usace.army.mil/cwms-data/levels'),
    ('a2d0956a-251c-4994-b8a8-3a240227ca4e', 'cwms-location', 'CWMS Location', 'https://cwms-data.usace.army.mil/cwms-data/location'),
    ('5d4a5d99-79e1-49b9-82aa-6f72708925e7', 'usgs-site', 'USGS Site', 'https://waterservices.usgs.gov/nwis/iv'),
    ('36dc9f8c-b18b-433c-b919-9c067739b6aa', 'usgs-timeseries', 'USGS Timeseries', 'https://waterservices.usgs.gov/nwis/iv'),
    ('42d3f5cb-8c5e-4857-a80d-202c0b86ed6c', 'nws-timeseries', 'NWS Timeseries', 'https://water.weather.gov'),
    ('f5854c47-7cc3-4bcb-9d13-83ed5ca31905', 'nws-level', 'NWS Level', 'https://water.weather.gov'),
    ('52f34db5-4129-41a3-812e-df63e7f9e715', 'nws-site', 'NWS Site', 'https://water.weather.gov');

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


INSERT INTO upload_status (id, name) VALUES
    ('b5d777fc-c46b-4a10-a488-1415e3d7849d', 'INITIATED'),
    ('969e00ad-2be1-4cf5-9f80-5c198e1e8450', 'SUCCESS'),
    ('020c8cda-913b-4c2d-8580-2834381bf885', 'FAIL');
