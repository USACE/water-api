#!/usr/bin/env python3

import requests
import json


states = ["AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DC", "DE", "FL", "GA", 
          "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", 
          "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", 
          "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", 
          "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"]


water_usgs_parameters_url = 'http://localhost/usgs_parameters'
r = requests.get(water_usgs_parameters_url)
parameter_codes = r.json()
water_usgs_parameters = []
for pc in parameter_codes:
    water_usgs_parameters.append(pc['code'])



for state in states:
    url = f'http://localhost/usgs_sites/state/{state.lower()}'
    #print(url)
    print("---------")
    print(state)
    print("---------")
    r = requests.get(url)
    #print(r.text)    

    sites = r.json()
    map = {}
    
    for s in sites:
        _p = {}
        for wp in water_usgs_parameters:
            # Check all USGS param codes in WATER API
            # and assigned True/False to sites params map
            if wp in s['parameter_codes']:
                _p[wp] = True
            else:
                _p[wp] = False        
        map[s['usgs_id']] = _p      
            

    # print('---------')
    # for k,v in map.items():
    #     print(k)
    #     if k == '03203600':
    #         print('\nHAAAAAA\n')
    #         print(k,v)

    insert_payload = []

    usgs_state_sites_url = f'https://waterservices.usgs.gov/nwis/iv/?format=json&stateCd={state.lower()}&siteType=LK,ST&siteStatus=active'
    req = requests.get(usgs_state_sites_url)
    usgs_state_sites = req.json()
    
    for i, state_site in enumerate(usgs_state_sites['value']['timeSeries']):
        # Note: Each site may be listed multiple times, once per parameter
        # print('-------')
        # print(state_site['sourceInfo']['siteName'])        
        usgs_id = state_site['sourceInfo']['siteCode'][0]['value']
        # print(usgs_id)
        param_code = state_site['variable']['variableCode'][0]['value']
        # print(param_code)
        if usgs_id in map.keys():

            if param_code in water_usgs_parameters and not map[usgs_id][param_code]:
                insert_payload.append({'usgs_id':usgs_id, 'usgs_parameter_codes':[param_code]})
                # print('payload inserted')
            # else:
            #     print(" -- PARAM EXISTS --")
        else:
            print('** USGS Site not in Water API')
        
    
    
    # print(usgs_state_sites['value']['timeSeries'][0]['variable']['variableCode'])
    for payload in insert_payload:
        print(payload)
    # Post payload for current state
    print("Sending payload...")
    r = requests.post(
    "http://localhost/usgs_sites/parameters?key=appkey",
    json=insert_payload,
    headers={"Content-Type": "application/json"},    
    )

    print(r.status_code)
    print(r.text)
    