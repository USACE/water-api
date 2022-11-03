-- This is required by seed data and needed before the repeatable migrations.
-- It can be updated/replaced by the repeatable Views migration.

--------------
-- V_LOCATION
--------------

CREATE OR REPLACE VIEW v_location AS (
    SELECT 
        l.id, 
        l.slug,
        l.code,
        l.geometry,
        sa.stusps AS state_abbrev,
        sa.name AS state,
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
        JOIN tiger_data.state_all sa ON sa.gid = l.state_id
);

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