#!/usr/bin/env python3

import requests
import uuid
import csv
from io import StringIO

"""
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
"""

states = {"TN": 30, "MN": 4, "WV": 1, "GA": 22, "SC": 28, "NC": 9, "AL": 23}

# states = ["AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DC", "DE", "FL", "GA",
#           "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
#           "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
#           "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
#           "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"]

site_types = {
    "LK": "896c45b0-e458-4621-840d-512f962427cb",
    "ST": "c34cb071-9163-4240-83fc-e0e691f61523",
    "ST-TS": "de46a4ab-b9ae-49d0-bac0-04fa07debc3d",
}

for state in states.keys():

    url = f"https://waterservices.usgs.gov/nwis/site/?format=rdb&stateCd={state}&period=P52W&siteType=LK,ST&siteStatus=all&hasDataTypeCd=iv,aw"
    # print(url)
    r = requests.get(url)
    # print(r.text)
    buff = StringIO(r.text)
    reader = csv.reader(buff, delimiter="\t")

    prepped_data = []
    keys = []
    result = []

    horizontal_datum = {"NAD83": 4269}
    vertical_datum = {"UNKNOWN": 0, "COE1912": 1, "NGVD29": 2, "NAVD88": 3}

    # Cleanup data before parse fields
    # Store header and actual data rows in new variable
    for idx, line in enumerate(reader):

        # only look at the header and data lines
        if line[0].strip() in ["agency_cd", "USGS"]:
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

    loc_sql = f"-- {state} sites\n"
    loc_sql += (
        "INSERT INTO location (id, datasource_id, slug, geometry, state_id) VALUES\n"
    )

    usgs_site_sql = ""
    usgs_site_sql += "INSERT INTO usgs_site (location_id, site_number, station_name, site_type_id) VALUES\n"
    last_line = len(result)

    for idx, line in enumerate(result):

        new_loc_id = uuid.uuid4()

        site_type = line["site_tp_cd"].strip()

        name = line["station_nm"].replace("'", "").strip()
        # elevation = (
        #     float(line["alt_va"].strip()) if line["alt_va"].strip() != "" else "NULL"
        # )
        # horizontal_datum_id = horizontal_datum[line["dec_coord_datum_cd"]]
        # huc = f"'{line['huc_cd'].strip()}'" if line["huc_cd"].strip() != "" else "NULL"
        # try:
        #     vertical_datum_id = vertical_datum[line["alt_datum_cd"]]
        # except:
        #     vertical_datum_id = vertical_datum["UNKNOWN"]

        loc_sql += f"('{new_loc_id}', '77dc8cf9-5804-434a-a53f-8b65c0358a6b', '{line['site_no'].strip()}', ST_GeomFromText('POINT({line['dec_long_va'].strip()} {line['dec_lat_va'].strip()})',4326), {states[state]})"

        usgs_site_sql += f"('{new_loc_id}', '{line['site_no'].strip()}', '{name}', '{site_types[site_type]}')"

        # loc_sql += f"('{line['site_no'].strip()}','{name}', ST_GeomFromText('POINT({line['dec_long_va'].strip()} {line['dec_lat_va'].strip()})',4326), "
        # loc_sql += f"{elevation}, {horizontal_datum_id},'{vertical_datum_id}', {huc}, '{state}')"

        # handle last set of values
        if idx + 1 == last_line:
            loc_sql += ";"
            usgs_site_sql += ";"
        else:
            loc_sql += ","
            usgs_site_sql += ","

        loc_sql += "\n"
        usgs_site_sql += "\n"

    print(loc_sql)

    print(usgs_site_sql)
