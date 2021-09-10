-- extent to polygon reference order - simple 4 point extents
-- xmin,ymax (top left), xmax ymax (top right), xmax ymin (bottom right), xmin ymin (bottom left), xmin ymax (top left again)
INSERT INTO watershed (id,slug,"name",geometry,office_id) VALUES
('03206ff6-fe91-426c-a5e9-4c651b06f9c6','eau-galla-river','Eau Galla River',ST_GeomFromText('POLYGON ((284000 2460000, 326000 2460000, 326000 2404000, 284000 2404000, 284000 2460000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
('048ce853-6642-4ac4-9fb2-81c01f67a85b','mississippi-river-headwaters','Mississippi River Headwaters',ST_GeomFromText('POLYGON ((24000 2760000, 254000 2760000, 254000 2402000, 24000 2402000, 24000 2760000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
('ad30f178-afc3-43b9-ba92-7bd139581217','red-river-north','Red River North',ST_GeomFromText('POLYGON ((-356000 2950000, 150000 2950000, 150000 2494000, -356000 2494000, -356000 2950000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
('c8bf6c6d-7f19-406a-a438-f2f876ce4815','souris-river','Souris River',ST_GeomFromText('POLYGON ((-708000 3100000, -178000 3100000, -178000 2736000, -708000 2736000, -708000 3100000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
('ced6ec9e-43b5-496e-a2b7-894af92c9b63','mississippi-river-navigation','Mississippi River Navigation',ST_GeomFromText('POLYGON ((48000 2646000, 564000 2646000, 564000 2204000, 48000 2204000, 48000 2646000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321'),
('f4219691-e498-46a3-ab0f-f2957bd09a10','minnesota-river','Minnesota River',ST_GeomFromText('POLYGON ((-112000 2602000, 234000 2602000, 234000 2244000, -112000 2244000, -112000 2602000))',5070),'2cf60156-f22a-418a-bc9f-a28960ad0321');
