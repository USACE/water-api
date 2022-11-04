DO $$

DECLARE 
    lrh_cwms_watershed_datasource_id uuid;
    lrn_cwms_watershed_datasource_id uuid;
    mvp_cwms_watershed_datasource_id uuid;
    sas_cwms_watershed_datasource_id uuid;
    lrl_cwms_watershed_datasource_id uuid;

BEGIN

    SELECT id into lrh_cwms_watershed_datasource_id FROM v_datasource  WHERE datatype = 'cwms-watershed' AND provider = 'lrh';
    SELECT id into lrn_cwms_watershed_datasource_id FROM v_datasource  WHERE datatype = 'cwms-watershed' AND provider = 'lrn';
    SELECT id into mvp_cwms_watershed_datasource_id FROM v_datasource  WHERE datatype = 'cwms-watershed' AND provider = 'mvp';
    SELECT id into sas_cwms_watershed_datasource_id FROM v_datasource  WHERE datatype = 'cwms-watershed' AND provider = 'sas';
    SELECT id into lrl_cwms_watershed_datasource_id FROM v_datasource  WHERE datatype = 'cwms-watershed' AND provider = 'lrl';


    INSERT INTO location (datasource_id, slug, code, geometry, state_id, attributes) VALUES
        -- Huntington District
        (lrh_cwms_watershed_datasource_id,'little-sandy-river','Little Sandy River',ST_GeomFromText('POLYGON ((1096000 1812000, 1158000 1812000, 1158000 1732000, 1096000 1732000, 1096000 1812000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'ohio-river','Ohio River',ST_GeomFromText('POLYGON ((1006000 1914000, 1206000 1914000, 1206000 1754000, 1006000 1754000, 1006000 1914000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'big-sandy-river','Big Sandy River',ST_GeomFromText('POLYGON ((1114000 1796000, 1288000 1796000, 1288000 1624000, 1114000 1624000, 1114000 1796000))',5070),NULL,'{}'), 
        (lrh_cwms_watershed_datasource_id,'twelvepole-river','Twelvepole River',ST_GeomFromText('POLYGON ((1152000 1796000, 1212000 1796000, 1212000 1728000, 1152000 1728000, 1152000 1796000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'hocking-river','Hocking River',ST_GeomFromText('POLYGON ((1112000 1960000, 1220000 1960000, 1220000 1878000, 1112000 1878000, 1112000 1960000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'guyandotte-river','Guyandotte River',ST_GeomFromText('POLYGON ((1166000 1814000, 1298000 1814000, 1298000 1692000, 1166000 1692000, 1166000 1814000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'scioto-river','Scioto River',ST_GeomFromText('POLYGON ((1004000 2056000, 1154000 2056000, 1154000 1810000, 1004000 1810000, 1004000 2056000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'kanawha-river','Kanawha River',ST_GeomFromText('POLYGON ((1182000 1870000, 1410000 1870000, 1410000 1544000, 1182000 1544000, 1182000 1870000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'muskingum-river','Muskingum River',ST_GeomFromText('POLYGON ((1098000 2110000, 1268000 2110000, 1268000 1904000, 1098000 1904000, 1098000 2110000))',5070),NULL,'{}'),
        (lrh_cwms_watershed_datasource_id,'little-kanawha-river','Little Kanawha River',ST_GeomFromText('POLYGON ((1164000 1970000, 1354000 1970000, 1354000 1824000, 1164000 1824000, 1164000 1970000))',5070),NULL,'{}'),
        -- Nashville District
        (lrn_cwms_watershed_datasource_id,'tennessee-river','Tennessee River',ST_GeomFromText('POLYGON ((640000 1678000, 1300000 1678000, 1300000 1268000, 640000 1268000, 640000 1678000))',5070),NULL,'{}'),
        (lrn_cwms_watershed_datasource_id,'cumberland-basin-river','Cumberland Basin River',ST_GeomFromText('POLYGON ((662000 1678000, 1172000 1678000, 1172000 1408000, 662000 1408000, 662000 1678000))',5070),NULL,'{}'),																												
        (lrn_cwms_watershed_datasource_id,'tennessee-cumberland-river', 'Tennessee & Cumberland River', ST_GeomFromText('POLYGON ((642000 1682000, 1300000 1682000, 1300000 1258000, 642000 1258000, 642000 1682000))',5070),NULL,'{}'),
        -- St. Paul District
        (mvp_cwms_watershed_datasource_id,'eau-galla-river','Eau Galla River',ST_GeomFromText('POLYGON ((284000 2460000, 326000 2460000, 326000 2404000, 284000 2404000, 284000 2460000))',5070),NULL,'{}'),
        (mvp_cwms_watershed_datasource_id,'mississippi-river-headwaters','Mississippi River Headwaters',ST_GeomFromText('POLYGON ((24000 2760000, 254000 2760000, 254000 2402000, 24000 2402000, 24000 2760000))',5070),NULL,'{}'),
        (mvp_cwms_watershed_datasource_id,'red-river-north','Red River North',ST_GeomFromText('POLYGON ((-356000 2950000, 150000 2950000, 150000 2494000, -356000 2494000, -356000 2950000))',5070),NULL,'{}'),
        (mvp_cwms_watershed_datasource_id,'souris-river','Souris River',ST_GeomFromText('POLYGON ((-708000 3100000, -178000 3100000, -178000 2736000, -708000 2736000, -708000 3100000))',5070),NULL,'{}'),
        (mvp_cwms_watershed_datasource_id,'mississippi-river-navigation','Mississippi River Navigation',ST_GeomFromText('POLYGON ((48000 2646000, 564000 2646000, 564000 2204000, 48000 2204000, 48000 2646000))',5070),NULL,'{}'),
        (mvp_cwms_watershed_datasource_id,'minnesota-river','Minnesota River',ST_GeomFromText('POLYGON ((-112000 2602000, 234000 2602000, 234000 2244000, -112000 2244000, -112000 2602000))',5070),NULL,'{}'),
        -- Savannah District
        (sas_cwms_watershed_datasource_id,'savannah-river-basin','Savannah River Basin',ST_GeomFromText('POLYGON ((1110000 1432000, 1432000 1432000, 1432000 1094000, 1110000 1094000, 1110000 1432000))',5070),NULL,'{}'),
        -- Louisville District
        (lrl_cwms_watershed_datasource_id,'middle-wabash-river','Middle Wabash River',ST_GeomFromText('POLYGON ((708000 1982000, 942000 1982000, 942000 1714000, 708000 1714000, 708000 1982000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'little-miami-river','Little Miami River',ST_GeomFromText('POLYGON ((980000 1952000, 1058000 1952000, 1058000 1824000, 980000 1824000, 980000 1952000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'licking-river','Licking River',ST_GeomFromText('POLYGON ((968000 1850000, 1150000 1850000, 1150000 1678000, 968000 1678000, 968000 1850000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'whitewater-river','Whitewater River',ST_GeomFromText('POLYGON ((900000 1958000, 968000 1958000, 968000 1844000, 900000 1844000, 900000 1958000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'green-river-1','Green River',ST_GeomFromText('POLYGON ((720000 1694000, 998000 1694000, 998000 1528000, 720000 1528000, 720000 1694000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'mill-creek-1','Mill Creek',ST_GeomFromText('POLYGON ((968000 1884000, 998000 1884000, 998000 1840000, 968000 1840000, 968000 1884000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'ohio-river-1','Ohio River',ST_GeomFromText('POLYGON ((584000 2004000, 1028000 2004000, 1028000 1566000, 584000 1566000, 584000 2004000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'great-miami-river','Great Miami River',ST_GeomFromText('POLYGON ((930000 2028000, 1056000 2028000, 1056000 1838000, 930000 1838000, 930000 2028000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'kentucky-river','Kentucky River',ST_GeomFromText('POLYGON ((924000 1818000, 1176000 1818000, 1176000 1606000, 924000 1606000, 924000 1818000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'little-wabash-river','Little Wabash River',ST_GeomFromText('POLYGON ((608000 1860000, 698000 1860000, 698000 1674000, 608000 1674000, 608000 1860000))',5070),NULL,'{}'),
        (lrl_cwms_watershed_datasource_id,'salt-river-1','Salt River',ST_GeomFromText('POLYGON ((866000 1766000, 980000 1766000, 980000 1646000, 866000 1646000, 866000 1766000))',5070),NULL,'{}');

END$$;
