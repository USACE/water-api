
-- allow location_id to be optional
ALTER TABLE visualization ALTER column location_id DROP NOT NULL;
-- add provider_id column
ALTER TABLE visualization ADD COLUMN provider_id UUID REFERENCES provider(id);

-- set the initial values (so they won't be NULL)
UPDATE visualization 
SET provider_id = (select office_id from location where id = visualization.location_id);

-- enforce the NOT NULL
ALTER TABLE visualization ALTER COLUMN provider_id SET NOT NULL;