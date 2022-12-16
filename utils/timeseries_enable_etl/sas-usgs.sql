DO $$

DECLARE

	sas_usgs_timeseries_id_array uuid[];
	timeseries_id uuid;

BEGIN

	sas_usgs_timeseries_id_array := ARRAY(
		SELECT id FROM water.v_timeseries WHERE DATATYPE = 'usgs-timeseries' AND KEY IN (
			'02176930.00065', '02176930.00060',
			'02177000.00065', '02177000.00060',
			'02178400.00065', '02178400.00060',
			'02181580.00065', '02181580.00060',
			'02186000.00065', '02186000.00060',
			'02187010.00062', '02187010.00045',
			'02187910.00065', '02187910.00060',
			'02188100.00062', '02188100.00045',
			'02188600.00065', '02188600.00060', '02188600.00045',
			'02192000.00065', '02192000.00060', '02192000.00045',
			'02192500.00065', '02192500.00060', '02192500.00045',
			'02193500.00065', '02193500.00060', '02193500.00045',
			'02193900.00062', '02193900.00045',
			'02195320.00065', '02195320.00060', '02195320.00045',
			'02195520.00065',
			'02196000.00065', '02196000.00060',
			'02196670.00065', '02196670.00060', '02196670.00045',
			'02196690.00065', '02196690.00060', '02196690.00045',
			'02196999.00065', '02196999.00060', '02196999.00045',
			'02197000.00065', '02197000.00060',
			'021973269.00065', '021973269.00060', '021973269.00045',
			'02197500.00065', '02197500.00060', '02197500.00045',
			'02198375.00065', '02198375.00060', '02198375.00045',
			'02198500.00065', '02198500.00060', '02198500.00045'
		)
	);
	
	
	

	--toggle the etl enabled
	UPDATE water.timeseries SET etl_values_enabled = TRUE 
	WHERE id = ANY(sas_usgs_timeseries_id_array);
	
	--create timeserie group	
	INSERT INTO water.timeseries_group (id, slug, name, provider_id) VALUES 
        ('9b232f39-57f6-414e-aa50-e0554e11bbd3', 'savannah-river', 'Savannah River Watershed Group', (SELECT id FROM water.provider WHERE slug = 'sas'));
       
    --insert group members
    FOREACH timeseries_id IN ARRAY sas_usgs_timeseries_id_array
    	LOOP
    		INSERT INTO water.timeseries_group_members (timeseries_group_id, timeseries_id) VALUES 
    		('9b232f39-57f6-414e-aa50-e0554e11bbd3', timeseries_id);
    END LOOP;
    	
    

END$$;