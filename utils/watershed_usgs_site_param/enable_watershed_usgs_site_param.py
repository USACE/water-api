# Python 3

import csv
import httpx  # need httpx and httpx[http2] installed

# WATER_API_ROOT = "https://develop-water-api.corps.cloud"
WATER_API_ROOT = "http://localhost"
WATERSHED_SLUG = "savannah-river-basin"
AUTH_TOKEN = ""


param_cwms_names = ["Stage", "Flow", "Elev-Pool", "Elevation"]
param_usgs_codes = ["00065", "00060", "00062", "00062"]

param_map = dict(zip(param_cwms_names, param_usgs_codes))


def get_water_usgs_params():
    r = httpx.get(f"{WATER_API_ROOT}/usgs/parameters")
    return r.json()


def post_watershed_site_param_config(ws_slug, site_number, param_code):

    url = f"{WATER_API_ROOT}/watersheds/{ws_slug}/site/{site_number}/parameter/{param_code}"

    if WATER_API_ROOT == "http://localhost":
        url = f"{url}?key=appkey"
        headers = None

    else:
        headers = {"Authorization": "Bearer " + AUTH_TOKEN}

    client = httpx.Client(http2=True)
    r = client.post(url=url, headers=headers)
    print(f"{url} -> {r.status_code}")
    if r.status_code not in ["201", "422"]:
        print(f"{r.text}")
    return {"status_code": r.status_code, "text": r.text}


# open file for reading
with open("sas_usgs_locations.csv") as csvDataFile:

    # read file as csv file
    csvReader = csv.reader(csvDataFile)

    success_config = []
    fail_config = []

    # for every row, print the row
    for row in csvReader:
        # print(row[0], row[5])
        site_number = row[0].strip()[1:][:-1]
        param_list = row[5]
        for param in param_list.split(","):
            # print(param)
            if param.strip() in param_cwms_names:
                # print(f"found {site_number} with a param of {param_map[param.strip()]}")
                req = post_watershed_site_param_config(
                    ws_slug=WATERSHED_SLUG,
                    site_number=site_number,
                    param_code=param_map[param.strip()],
                )
                if req["status_code"] == 201:
                    success_config.append([site_number, param_map[param.strip()]])
                else:
                    fail_config.append([site_number, param_map[param.strip()]])

    print(f"Success Count: {len(success_config)}")
    print(f"Fail Count: {len(fail_config)}")
