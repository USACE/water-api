--------------
-- V_DATASOURCE
--------------

CREATE OR REPLACE VIEW v_datasource AS (
    SELECT
        d.id,
        dt.slug as datatype,
        dt."name" datatype_name,
        dt.uri,
        p.slug AS provider,
        p."name" AS provider_name
        FROM datasource d
        JOIN "datatype" dt ON dt.id = d.datatype_id 
        JOIN provider p ON p.id = d.provider_id 
);

--------------
-- V_LOCATION
--------------

CREATE OR REPLACE VIEW v_location AS (
    SELECT 
        l.id, 
        l.slug,
        l.code,
        l.geometry,
        sa.stusps AS state,
        sa.name AS state_name,
        l.create_date,
        l.update_date,
        l."attributes",
        p."name" AS provider_name,
        p.slug AS provider,
        dt.slug AS "datatype",
        dt."name" AS datatype_name
        FROM "location" l
        JOIN datasource d ON d.id = l.datasource_id
        JOIN "datatype" dt ON dt.id = d.datatype_id
        JOIN provider p ON p.id = d.provider_id 
        LEFT JOIN tiger_data.state_all sa ON sa.gid = l.state_id
);


---------------
-- V_TIMESERIES
---------------

CREATE OR REPLACE VIEW v_timeseries AS (
    SELECT t.id              AS id,
           p1.slug  		 AS provider,
		   p1.name           AS provider_name,
		   dt1.slug 		 AS datatype,
           dt1.name          AS datatype_name,
		   t.datasource_key  AS key,
		   CASE
               WHEN t.latest_time IS NULL OR t.latest_value IS NULL THEN NULL
               ELSE json_build_array(t.latest_time, t.latest_value)::json
           END AS latest_value,
           json_build_object(
               'slug'    ,   l.slug,
               'provider',  p2.slug,
               'datatype', dt2.slug,
               'code'    ,   l.code
		   )                            AS location,
           t.etl_values_enabled         AS etl_values_enabled
    FROM timeseries t
    JOIN datasource   ds1 ON ds1.id =   t.datasource_id  -- timeseries' datasource
    JOIN provider      p1 ON  p1.id = ds1.provider_id    -- timeseries' provider
    JOIN datatype     dt1 ON dt1.id = ds1.datatype_id    -- timeseries' datatype
    JOIN location       l ON   l.id =   t.location_id
    JOIN datasource   ds2 ON ds2.id =   l.datasource_id  -- location's datasource
    JOIN provider      p2 ON  p2.id = ds2.provider_id    -- location's provider
    JOIN datatype     dt2 ON dt2.id = ds2.datatype_id    -- location's datatype
);

--------------
-- V_WATERSHED
--------------
-- TODO; rethink or refactor this view
-- This should rebuild after being deleted in 1.0.7
CREATE OR REPLACE VIEW v_watershed AS (
    SELECT l.id       AS id,
           l.slug     AS slug,
           l.code     AS name,
           l.geometry AS geometry,
           p.id       AS provider_id,
           p.slug     AS provider_slug
	FROM   location l
    JOIN datasource d on d.id = l.datasource_id
        AND d.datatype_id = (SELECT id FROM datatype WHERE slug = 'cwms-watershed')
    JOIN provider p ON p.id = d.provider_id
);
