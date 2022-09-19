DROP VIEW IF EXISTS v_watershed;

ALTER TABLE watershed
    ALTER COLUMN geometry
        TYPE geometry(Geometry, 4326)
        USING ST_Transform(geometry, 4326);


ALTER TABLE watershed
    ALTER COLUMN geometry
        SET DEFAULT ST_GeomFromText('POLYGON ((0 0, 0 0, 0 0, 0 0, 0 0))',4326);


-- The view should get replaced in R__01_Views.sql
-- Need to make minor change (just a comment) to force
-- a new file hash.
