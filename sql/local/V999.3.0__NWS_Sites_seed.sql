DO $$

DECLARE 
nws_site_ds_id uuid;

BEGIN

    SELECT d.id into nws_site_ds_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'nws-site' AND p.slug = 'noaa-nws';

INSERT INTO location (id, datasource_id, slug, geometry, state_id) VALUES
('126510d7-6394-4589-891a-320818efd978', nws_site_ds_id, 'kanw2', ST_GeomFromText('POINT(-0.0 0.0)',4326), 1);

INSERT INTO nws_site (location_id, name, nws_li) VALUES
('126510d7-6394-4589-891a-320818efd978', 'Kanawha Falls at Kanawha River', 'kanw2');


END$$;