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

        INSERT into datasource(provider_id, datasource_type_id) 
        VALUES (
            (SELECT id from provider where lower(slug) = lower(cwms_provider)), 
            (SELECT id from datasource_type where lower(slug) = cwms_data_type)
        );
        
        END LOOP;

    END LOOP;

END$$;




INSERT into datasource(id, provider_id, datasource_type_id) VALUES
-- USGS
('77dc8cf9-5804-434a-a53f-8b65c0358a6b', (SELECT id from provider where slug = 'usgs'), (SELECT id from datasource_type where slug = 'usgs-timeseries')),
-- NWS
('a59ffe5f-6614-4679-a387-204013aa8de3', (SELECT id from provider where slug = 'nws'), (SELECT id from datasource_type where slug = 'nws-timeseries')),
('9aafe174-dc17-44c8-813e-714474b44e04', (SELECT id from provider where slug = 'nws'), (SELECT id from datasource_type where slug = 'nws-level'));