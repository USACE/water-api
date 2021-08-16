#!/usr/bin/env python3

import requests
import json
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

# states = ['TN', 'MN', 'WV', 'GA', 'SC', 'NC', 'AL']

states = ["AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DC", "DE", "FL", "GA", 
          "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", 
          "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", 
          "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", 
          "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"]

# states = ['MA']

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


    # for line in result:
    #     print('-'*10)
    #     for k,v in line.items():
    #         print(k, '=>', v)

    state_sites = []    

    # Each line (after cleanup) in the results from the USGS API call represents a site
    for line in result:

        site = {}
        site['usgs_id'] = line['site_no'].strip()
        site['name'] = line['station_nm'].replace("'", "").strip()
        site['state_abbrev'] = state
        site['elevation'] = float(line['alt_va'].strip()) if line['alt_va'].strip() != '' else None
        try:
            site['horizontal_datum_id'] = horizontal_datum[line['dec_coord_datum_cd']]
        except:
            site['horizontal_datum_id'] = 4269
        site['huc'] = f"{line['huc_cd'].strip()}" if line['huc_cd'].strip() != '' else None
        try:
            site['vertical_datum_id'] = vertical_datum[line['alt_datum_cd']]
        except:
            site['vertical_datum_id'] = vertical_datum['UNKNOWN']
        geom = {}
        geom['type'] = 'Point'
        try:
            geom['coordinates'] = [float(line['dec_long_va'].strip()), float(line['dec_lat_va'].strip())]
        except:
            geom['coordinates'] = [0,0]
        site['geometry'] = geom

  
        state_sites.append(site)
    
    # print(state_sites)
    
    
    r = requests.post(
    "http://localhost/sync/usgs_sites?key=appkey",
    json=state_sites,
    headers={"Content-Type": "application/json"},    
    )
    # print(json.dumps(state_sites, indent=4))
    print(r.status_code)
    print(r.text)
