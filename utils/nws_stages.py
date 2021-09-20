import requests
import csv
from io import StringIO

# Get USGS sites from API
url = 'http://localhost/usgs/sites'
r = requests.get(url)
sites = r.json()

usgs_sites = []
for site in sites:
    usgs_sites.append(site['site_number'])

# Get Existing stages from API
url = 'http://localhost/nws/stages'
r = requests.get(url)
stages = r.json()

existing_stages = {}
for site_stages in stages:
    existing_stages[site_stages['nwsid']] = site_stages


# Get states from API
url = 'http://localhost/states'
r = requests.get(url)
api_states = r.json()

states = []
for state in api_states:
    states.append(state['abbreviation'])




url = "https://water.weather.gov/monitor/ahpsreport.php?type=csv"
# print(url)
r = requests.get(url)
# print(r.text)
buff = StringIO(r.text)
reader = csv.reader(buff, delimiter=',')

keys = []
result= []

for idx, line in enumerate(reader):
    if idx == 0:
        # this is the header
        keys = line
        # print(keys)
    else:
        # Build each line (object) by setting the keys and values
        _line = {}
        for i, k in enumerate(keys):
            _line[k] = line[i].strip()
        result.append(_line)

# for line in result:
#         print('-'*10)
#         for k,v in line.items():
#             print(k, '=>', v)


for wstate in states:

    sql = '-- '+wstate
    sql += '\nINSERT INTO nws_stages (nwsid, usgs_site_number, name, action, flood, moderate_flood, major_flood) VALUES\n'
    last_line = len(result)
    for idx, line in enumerate(result):

        if line['state'].strip() == wstate.lower():
        
            # Entry must have:
            # 1) USGSID must be in API sites
            # 2) NWSID = 5 chars
            # 3) a proper USGS ID 
            # 4) stage must have a unit of 'FT
            # 5) all stages cannot be 0
            if line['usgs id'].strip() in usgs_sites and len(line['nws shef id'].strip()) == 5 and line['usgs id'].strip() != '' \
                and len(line['usgs id'].strip()) >= 8 and line['usgs id'].strip().isnumeric() and line['flood stage unit'].strip() == 'FT' \
                and float(line['action stage'].strip()) != 0 and float(line['flood stage'].strip()) != 0 \
                and float(line['moderate flood stage'].strip()) != 0 and float(line['major flood stage'].strip()) != 0:

                if line['proximity'].strip() in ['at', 'near', 'below', 'near']:
                    # ex: Name above xyz lake
                    name = f"{line['location name'].strip()} {line['proximity'].strip()} {line['river/water-body name'].strip()}"
                else:
                    name = line['location name'].strip()

                nwsid                = line['nws shef id'].strip()
                usgs_id              = line['usgs id'].strip()
                action         = float(line['action stage'].strip())
                flood          = float(line['flood stage'].strip())
                moderate_flood = float(line['moderate flood stage'].strip())
                major_flood    = float(line['major flood stage'].strip())

                
                sql += f"('{nwsid}', '{usgs_id}', '{name}', {action}, {flood}, {moderate_flood}, {major_flood}),\n"

                # API Post Payload
                if nwsid in existing_stages.keys():
                    # site stages already in database, test for changes, if changes use put
                    # print('NWS site already in DB')
                    continue
                else:
                    # new site stages record, use POST                
                    
                    payload = {}
                    payload['nwsid'] = nwsid
                    payload['usgs_site_number'] = usgs_id
                    payload['name'] = name
                    payload['action'] = action
                    payload['flood'] = flood
                    payload['moderate_flood'] = moderate_flood
                    payload['major_flood'] = major_flood


                    r = requests.post(
                    "http://localhost/nws/stages?key=appkey",
                    json=payload,
                    headers={"Content-Type": "application/json"},    
                    )
                    if r.status_code != 201:
                        print(f'Unable to post {nwsid}')
                        print(r.text)
            
            
    # semi-colon on last line
    sql = sql[0:-2]+';'
            
    print(sql)
