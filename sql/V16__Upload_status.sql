-- upload_status definition
CREATE TABLE IF NOT EXISTS upload_status (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL
);

-- grant read
GRANT SELECT ON
    upload_status
TO water_reader;

-- grant write
GRANT INSERT,UPDATE,DELETE ON
    upload_status
TO water_writer;

INSERT INTO upload_status (id, name) VALUES
    ('b5d777fc-c46b-4a10-a488-1415e3d7849d', 'INITIATED'),
    ('969e00ad-2be1-4cf5-9f80-5c198e1e8450', 'SUCCESS'),
    ('020c8cda-913b-4c2d-8580-2834381bf885', 'FAIL');