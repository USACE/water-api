-- HUC2 Watersheds
-- QGIS Export to PostGIS Format Using File:
BEGIN;
CREATE TABLE "a2w_cwms"."usgs_huc2" ("ogc_fid" SERIAL,CONSTRAINT "usgs_huc2_pk" PRIMARY KEY ("ogc_fid"));
SELECT AddGeometryColumn('a2w_cwms','usgs_huc2','wkb_geometry',4326,'MULTIPOLYGON',2);
CREATE INDEX "usgs_huc2_wkb_geometry_geom_idx" ON "a2w_cwms"."usgs_huc2" USING GIST ("wkb_geometry");
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "objectid" NUMERIC(20,0);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "tnmid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "metasourceid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "sourcedatadesc" VARCHAR(100);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "sourceoriginator" VARCHAR(130);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "sourcefeatureid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "loaddate" timestamp with time zone;
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "areasqkm" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "areaacres" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "referencegnis_ids" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "name" VARCHAR(120);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "states" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "huc2" VARCHAR(2);
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "shape_length" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc2" ADD COLUMN "shape_area" FLOAT8;
COMMIT;

-- SIMPLIFIED HUC2 WATERSHEDS WITH SAME DATA SCHEMA
CREATE TABLE a2w_cwms.usgs_huc2_simple AS TABLE a2w_cwms.usgs_huc2;

-- HUC4 Watersheds
-- QGIS Export to PostGIS Format Using File:
BEGIN;
CREATE TABLE "a2w_cwms"."usgs_huc4" ( "ogc_fid" SERIAL, CONSTRAINT "usgs_huc4_pk" PRIMARY KEY ("ogc_fid") );
SELECT AddGeometryColumn('a2w_cwms','usgs_huc4','wkb_geometry',4326,'MULTIPOLYGON',2);
CREATE INDEX "usgs_huc4_wkb_geometry_geom_idx" ON "a2w_cwms"."usgs_huc4" USING GIST ("wkb_geometry");
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "objectid" NUMERIC(20,0);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "tnmid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "metasourceid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "sourcedatadesc" VARCHAR(100);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "sourceoriginator" VARCHAR(130);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "sourcefeatureid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "loaddate" timestamp with time zone;
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "areasqkm" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "areaacres" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "referencegnis_ids" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "name" VARCHAR(120);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "states" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "huc4" VARCHAR(4);
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "shape_length" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc4" ADD COLUMN "shape_area" FLOAT8;
COMMIT;

-- SIMPLIFIED HUC4 WATERSHEDS WITH SAME DATA SCHEMA
CREATE TABLE a2w_cwms.usgs_huc4_simple AS TABLE a2w_cwms.usgs_huc4;

-- HUC6 Watersheds
-- QGIS Export to PostGIS Format Using File:
BEGIN;
CREATE TABLE "a2w_cwms"."usgs_huc6" ( "ogc_fid" SERIAL, CONSTRAINT "usgs_huc6_pk" PRIMARY KEY ("ogc_fid") );
SELECT AddGeometryColumn('a2w_cwms','usgs_huc6','wkb_geometry',4326,'MULTIPOLYGON',2);
CREATE INDEX "usgs_huc6_wkb_geometry_geom_idx" ON "a2w_cwms"."usgs_huc6" USING GIST ("wkb_geometry");
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "objectid" NUMERIC(20,0);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "tnmid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "metasourceid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "sourcedatadesc" VARCHAR(100);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "sourceoriginator" VARCHAR(130);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "sourcefeatureid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "loaddate" timestamp with time zone;
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "areasqkm" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "areaacres" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "referencegnis_ids" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "name" VARCHAR(120);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "states" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "huc6" VARCHAR(6);
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "shape_length" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc6" ADD COLUMN "shape_area" FLOAT8;
COMMIT;

-- SIMPLIFIED HUC6 WATERSHEDS WITH SAME DATA SCHEMA
CREATE TABLE a2w_cwms.usgs_huc6_simple AS TABLE a2w_cwms.usgs_huc6;


-- HUC8 Watersheds
-- QGIS Export to PostGIS Format Using File:
BEGIN;
CREATE TABLE "a2w_cwms"."usgs_huc8" ( "ogc_fid" SERIAL, CONSTRAINT "usgs_huc8_pk" PRIMARY KEY ("ogc_fid") );
SELECT AddGeometryColumn('a2w_cwms','usgs_huc8','wkb_geometry',4326,'MULTIPOLYGON',2);
CREATE INDEX "usgs_huc8_wkb_geometry_geom_idx" ON "a2w_cwms"."usgs_huc8" USING GIST ("wkb_geometry");
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "objectid" NUMERIC(20,0);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "tnmid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "metasourceid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "sourcedatadesc" VARCHAR(100);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "sourceoriginator" VARCHAR(130);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "sourcefeatureid" VARCHAR(40);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "loaddate" timestamp with time zone;
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "areasqkm" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "areaacres" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "referencegnis_ids" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "name" VARCHAR(120);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "states" VARCHAR(50);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "huc8" VARCHAR(8);
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "shape_length" FLOAT8;
ALTER TABLE "a2w_cwms"."usgs_huc8" ADD COLUMN "shape_area" FLOAT8;
COMMIT;

-- SIMPLIFIED HUC8 WATERSHEDS WITH SAME DATA SCHEMA
CREATE TABLE a2w_cwms.usgs_huc8_simple AS TABLE a2w_cwms.usgs_huc8;

GRANT SELECT ON
    usgs_huc2,
    usgs_huc2_simple,
    usgs_huc4,
    usgs_huc4_simple,
    usgs_huc6,
    usgs_huc6_simple,
    usgs_huc8,
    usgs_huc8_simple
TO water_reader;
