#!/usr/bin/env python3

import requests
import csv
from io import StringIO
'''
Header fields:
#  agency_cd       -- Agency
#  site_no         -- Site identification number
#  station_nm      -- Site name
#  site_tp_cd      -- Site type
#  dec_lat_va      -- Decimal latitude
#  dec_long_va     -- Decimal longitude
#  coord_acy_cd    -- Latitude-longitude accuracy
#  dec_coord_datum_cd -- Decimal Latitude-longitude datum
#  alt_va          -- Altitude of Gage/land surface
#  alt_acy_va      -- Altitude accuracy
#  alt_datum_cd    -- Altitude datum
#  huc_cd          -- Hydrologic unit code

https://waterservices.usgs.gov/rest/Site-Service.html#Understanding
'''

states = ['TN', 'MN', 'WV', 'GA', 'SC', 'NC', 'AL']

# states = ["AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DC", "DE", "FL", "GA", 
#           "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", 
#           "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", 
#           "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", 
#           "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"]

for state in states:

    url = f"https://waterservices.usgs.gov/nwis/site/?format=rdb&stateCd={state}&period=P52W&siteType=LK,ST&siteStatus=all&hasDataTypeCd=iv,aw"
    # print(url)
    r = requests.get(url)
    # print(r.text)
    buff = StringIO(r.text)
    reader = csv.reader(buff, delimiter='\t')

    prepped_data = []
    keys = []
    result= []

    horizontal_datum = {
        'NAD83': 4269
    }
    vertical_datum = {
        'UNKNOWN': 0,
        'COE1912': 1,
        'NGVD29': 2,
        'NAVD88': 3        
    }


    # Cleanup data before parse fields
    # Store header and actual data rows in new variable
    for idx, line in enumerate(reader):
        
        # only look at the header and data lines
        if line[0].strip() in ['agency_cd', 'USGS']:
            prepped_data.append(line)

    for idx, line in enumerate(prepped_data):
            
        # print(idx)
        if idx == 0:
            # this is the header
            keys = line
        else:
            # Build each line (object) by setting the keys and values
            _line = {}
            for i, k in enumerate(keys):
                _line[k] = line[i].strip()
            result.append(_line)

    # print(result)

    # for line in result:
    #     print('-'*10)
    #     for k,v in line.items():
    #         print(k, '=>', v)


    sql = f"-- {state} sites\n"
    sql += 'INSERT INTO usgs_site (site_number, name, geometry, elevation, horizontal_datum_id, vertical_datum_id, huc, state_abbrev) VALUES\n'
    last_line = len(result)
    for idx, line in enumerate(result):

        name = line['station_nm'].replace("'", "").strip()
        elevation = float(line['alt_va'].strip()) if line['alt_va'].strip() != '' else 'NULL'
        horizontal_datum_id = horizontal_datum[line['dec_coord_datum_cd']]
        huc = f"'{line['huc_cd'].strip()}'" if line['huc_cd'].strip() != '' else 'NULL'
        try:
            vertical_datum_id = vertical_datum[line['alt_datum_cd']]
        except:
            vertical_datum_id = vertical_datum['UNKNOWN']

        sql += f"('{line['site_no'].strip()}','{name}', ST_GeomFromText('POINT({line['dec_long_va'].strip()} {line['dec_lat_va'].strip()})',4326), "
        sql += f"{elevation}, {horizontal_datum_id},'{vertical_datum_id}', {huc}, '{state}')"

        # handle last set of values
        if idx+1 == last_line:
            sql += ';'
        else:
            sql += ','

        sql += '\n'

    print(sql)