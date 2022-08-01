DROP VIEW IF EXISTS v_watershed;

ALTER TABLE watershed
    ALTER COLUMN geometry
        TYPE geometry(Geometry, 4326)
        USING ST_Transform(geometry, 4326);


ALTER TABLE watershed
    ALTER COLUMN geometry
        SET DEFAULT ST_GeomFromText('POLYGON ((0 0, 0 0, 0 0, 0 0, 0 0))',4326);


CREATE OR REPLACE VIEW v_watershed AS (
    SELECT w.id,
           w.slug,
           w.name,
           w.geometry AS geometry,
           w.office_id,
           f.symbol AS office_symbol
	FROM   watershed w
    LEFT JOIN office f ON w.office_id = f.id
	WHERE NOT w.deleted
);
