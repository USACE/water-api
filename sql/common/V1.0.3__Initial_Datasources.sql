DO $$

DECLARE
cwms_provider_array varchar[] := '{LRB, LRC, LRD, LRE, LRH, LRL, LRN, LRP, MVD, MVK, MVM, MVN, MVP, MVR, MVS, NAB, NAD, NAE, NAN, NAO, NAP, NWD, '
                                 'NWK, NWO, NWP, NWS, NWW, POA, POH, SAC, SAD, SAJ, SAM, SAS, SAW, SPA, SPD, SPK, SPL, SPN, SWD, SWF, SWL, SWT}';
cwms_provider varchar;

cwms_data_type_array varchar[] := '{cwms-timeseries, cwms-level, cwms-location}';
cwms_data_type varchar;

BEGIN
FOREACH cwms_provider IN ARRAY cwms_provider_array
	LOOP

    -- loop over each cwms data type
    FOREACH cwms_data_type IN ARRAY cwms_data_type_array
    LOOP

        INSERT into datasource(provider_id, datatype_id) 
        VALUES (
            (SELECT id from provider where lower(slug) = lower(cwms_provider)), 
            (SELECT id from datatype where lower(slug) = cwms_data_type)
        );
        
        END LOOP;

    END LOOP;

END$$;




INSERT into datasource(id, provider_id, datatype_id) VALUES
-- USGS
('2988a8d9-19a2-4c54-8594-f546a32fe43c', (SELECT id from provider where slug = 'usgs'), (SELECT id from datatype where slug = 'usgs-site')),
('77dc8cf9-5804-434a-a53f-8b65c0358a6b', (SELECT id from provider where slug = 'usgs'), (SELECT id from datatype where slug = 'usgs-timeseries')),
-- NWS
('a59ffe5f-6614-4679-a387-204013aa8de3', (SELECT id from provider where slug = 'noaa-nws'), (SELECT id from datatype where slug = 'nws-timeseries')),
('9aafe174-dc17-44c8-813e-714474b44e04', (SELECT id from provider where slug = 'noaa-nws'), (SELECT id from datatype where slug = 'nws-level')),
('65195574-602c-4540-8ba2-60198ca8e17e', (SELECT id from provider where slug = 'noaa-nws'), (SELECT id from datatype where slug = 'nws-site'));