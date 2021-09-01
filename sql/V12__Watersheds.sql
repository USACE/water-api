-- watershed
CREATE TABLE IF NOT EXISTS watershed (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    slug VARCHAR UNIQUE NOT NULL,
    name VARCHAR,
    geometry geometry NOT NULL DEFAULT ST_GeomFromText('POLYGON ((0 0, 0 0, 0 0, 0 0, 0 0))',5070),
    office_id UUID NOT NULL REFERENCES office(id),
	deleted boolean NOT NULL DEFAULT false
);

-- Grant read
GRANT SELECT ON
    watershed
TO water_reader;

-- Grant write
GRANT INSERT,UPDATE,DELETE ON
    watershed
TO water_writer;

-- extent to polygon reference order - simple 4 point extents
-- xmin,ymax (top left), xmax ymax (top right), xmax ymin (bottom right), xmin ymin (bottom left), xmin ymax (top left again)
INSERT INTO watershed (id,slug,"name",geometry,office_id) VALUES	 
    ('0f065e6a-3380-4ac3-b576-89fae7774b9f','little-sandy-river','Little Sandy River',ST_GeomFromText('POLYGON ((1096000 1812000, 1158000 1812000, 1158000 1732000, 1096000 1732000, 1096000 1812000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),    
    ('1a629fac-82c9-4b3e-b7fc-6a891d944140','ohio-river','Ohio River',ST_GeomFromText('POLYGON ((1006000 1914000, 1206000 1914000, 1206000 1754000, 1006000 1754000, 1006000 1914000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	
    ('3e322a11-b76b-4710-8f9a-b7884cd8ae77','big-sandy-river','Big Sandy River',ST_GeomFromText('POLYGON ((1114000 1796000, 1288000 1796000, 1288000 1624000, 1114000 1624000, 1114000 1796000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('4d3083d1-101c-4b76-9311-1154917ffbf1','twelvepole-river','Twelvepole River',ST_GeomFromText('POLYGON ((1152000 1796000, 1212000 1796000, 1212000 1728000, 1152000 1728000, 1152000 1796000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('5024720e-02f6-4577-a09c-ff1ff5a28223','hocking-river','Hocking River',ST_GeomFromText('POLYGON ((1112000 1960000, 1220000 1960000, 1220000 1878000, 1112000 1878000, 1112000 1960000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('50372dbc-f254-4584-8345-1c3613d2a102','guyandotte-river','Guyandotte River',ST_GeomFromText('POLYGON ((1166000 1814000, 1298000 1814000, 1298000 1692000, 1166000 1692000, 1166000 1814000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('5758d0dc-c8bf-4e37-a5e7-44ff3f4b8677','scioto-river','Scioto River',ST_GeomFromText('POLYGON ((1004000 2056000, 1154000 2056000, 1154000 1810000, 1004000 1810000, 1004000 2056000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('65a93467-c9b4-4166-acb6-58e8ec06ed3b','kanawha-river','Kanawha River',ST_GeomFromText('POLYGON ((1182000 1870000, 1410000 1870000, 1410000 1544000, 1182000 1544000, 1182000 1870000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),	 
    ('7c6dd902-fbc5-43e4-9bbf-351963f5723d','muskingum-river','Muskingum River',ST_GeomFromText('POLYGON ((1098000 2110000, 1268000 2110000, 1268000 1904000, 1098000 1904000, 1098000 2110000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
    ('cf193b4e-61c3-4e4d-9503-2935a82aed96','little-kanawha-river','Little Kanawha River',ST_GeomFromText('POLYGON ((1164000 1970000, 1354000 1970000, 1354000 1824000, 1164000 1824000, 1164000 1970000))',5070),'2f160ba7-fd5f-4716-8ced-4a29f75065a6'),
	('c54eab5b-1020-476b-a5f8-56d77802d9bf','tennessee-river','Tennessee River',ST_GeomFromText('POLYGON ((640000 1678000, 1300000 1678000, 1300000 1268000, 640000 1268000, 640000 1678000))',5070),'552e59f7-c0cc-4689-8a4d-e791c028430a'),	 
    ('c785f4de-ab17-444b-b6e6-6f1ad16676e8','cumberland-basin-river','Cumberland Basin River',ST_GeomFromText('POLYGON ((662000 1678000, 1172000 1678000, 1172000 1408000, 662000 1408000, 662000 1678000))',5070),'552e59f7-c0cc-4689-8a4d-e791c028430a'),	 																																	
	('feda585b-1ba0-4b19-92ed-7195154b8052','tennessee-cumberland-river', 'Tennessee & Cumberland River', ST_GeomFromText('POLYGON ((642000 1682000, 1300000 1682000, 1300000 1258000, 642000 1258000, 642000 1682000))',5070), '552e59f7-c0cc-4689-8a4d-e791c028430a');	 


-- -----
-- VIEWS
-- -----

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

-- Grant read
GRANT SELECT ON
    v_watershed
TO water_reader;