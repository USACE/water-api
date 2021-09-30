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

-- grant read
GRANT SELECT ON
    watershed_shapefile_uploads
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    watershed_shapefile_uploads
TO water_writer;

-- seed data
-- INSERT INTO watershed_shapefile_uploads (watershed_id, file, file_size) VALUES
--     ('c785f4de-ab17-444b-b6e6-6f1ad16676e8','water/watersheds/cumberland-basin-river/watersheds.zip', 1);