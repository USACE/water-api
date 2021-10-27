# Python 3

import csv
import httpx  # need httpx and httpx[http2] installed

WATER_API_ROOT = "https://develop-water-api.corps.cloud"
# WATER_API_ROOT = "http://localhost"
WATERSHED_SLUG = "savannah-river-basin"
AUTH_TOKEN = "eyJhbGciOiJSUzUxMiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJiWWFWUk0tMVZqdS1kdjdjRGdJOWZyZDVUbWRZdUVOUG1kaFdDWmVNU1pjIn0.eyJleHAiOjE2MzUzNjUxMDUsImlhdCI6MTYzNTM2NDgwNSwianRpIjoiMGNlNDlmYTEtZTY5NC00OTc5LWJiZjktYTI2ODJjN2IwZGMyIiwiaXNzIjoiaHR0cHM6Ly9hdXRoLmNvcnBzLmNsb3VkL2F1dGgvcmVhbG1zL3dhdGVyIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6ImQzZjg4Y2JkLTUzMDAtNDMyYS05NzE4LWI0NjA2OTAwMDYxOSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImEydyIsInNlc3Npb25fc3RhdGUiOiIwZmNmZTg2OS01MzZhLTQxZWQtYjAwMC05MzJlZDRkNWJhNGYiLCJhY3IiOiIxIiwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImRlZmF1bHQtcm9sZXMtd2F0ZXIiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYTJ3Ijp7InJvbGVzIjpbIkFQUExJQ0FUSU9OLkFETUlOIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIiwicHJlZmVycmVkX3VzZXJuYW1lIjoic2NhcmJlcnJ5LnJhbmRhbC5hIiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIn0.h0Cdk-il8DYNGv_T3LzY5gfXK_nhFITqC0eA7JDqrXkK1-NlP_MAO4EjqZdeQCPJqoc0cITsHCm2orR79c82kV8pTiFjqajwUOWOOxITaaI5yKSgCObhSvM-dhQ37LNtFKTp5N-As-yo0c41lb0OAvKim9E7Jy_2cXAF7hzed0KDmbquABTaAsZ4tTAmh9YK6Ycku7abAzhpLXXzoTzGY11bSN1aF-5_TaKkcLIXzu5SPXEpKo_H0Yg9vZuiYKxGfq9FS7KLbupO8wtOcumYFmiGBtJsF5Di6ymS7Sj9peIsyWbY9s-dIrbI4z6rrFdbYhU-Y3P3F3vUWmEdqBrxuaoSsWvOZhsKGtLES_iisB-_hb0wpCPBAyK2lKw5dv9PDTQqwbw93H0zVRzjLADyDJJe9pIh6DfveXGAhYtJGEJEo-7ELVcjgyvzR6_zrc-RL3IYLmf7jSwzfyrcyfb-iMUJ1i0MZ7N5d9L-0xexhT0jHfcrAnDqkbWHFCv_wB-eUaJud1LZmSUkMtxdv48pYyySA83xyhGvl_F3DvHPEVTBtSj3P8hISR2eJJbVkliM5wzcx3qOowYIDxEL6VAF997oJFD5wrxZNSYJ4BUUHrtPNQBJK9Br9RDYfCttMgntB-ujTezAiAuTcGYoMCtwbnjQPQKuzlCTeJzGbA6TfnM"


param_cwms_names = ["Stage", "Flow", "Elev-Pool", "Elevation"]
param_usgs_codes = ["00065", "00060", "00062", "00062"]

param_map = dict(zip(param_cwms_names, param_usgs_codes))


def get_water_usgs_params():
    r = httpx.get(f"{WATER_API_ROOT}/usgs/parameters")
    return r.json()


def post_watershed_site_param_config(ws_slug, site_number, param_code):

    url = f"{WATER_API_ROOT}/watersheds/{ws_slug}/site/{site_number}/parameter/{param_code}"

    if WATER_API_ROOT == "localhost":
        url = f"{url}?key=appkey"

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
