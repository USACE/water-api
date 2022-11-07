import requests
import re
import uuid
import json
from slugify import slugify  # lib is python-slugify


def apostrophe_doubleup(input):
    return input.replace("'", "''")


def quote(input):
    """Wrap in single quotes and double-up any apostrophes ' --> '' for valid SQL Inserts"""
    return f"'{input}'"


def nextUniqueSlug(base_text, slug_map):
    """Returns a unique slug given a base string and map of existing slugs (to ensure uniqueness)
    A slug is a string with spaces converted to "-" and special characters removed, useful for display in URLs
    """

    baseSlug = slugify(base_text)
    # If Slug Already Unique, Use It
    if baseSlug not in slug_map.keys():
        return baseSlug
    # If Slug is Not Unique, Iterate and append "-1", "-2", "-3" until slug is unique
    _slug, suffix = baseSlug, 1
    while _slug in slug_map.keys():
        _slug = f"{baseSlug}-{suffix}"
        suffix += 1

    return _slug


def radar_fix_json(text):
    def replacement(matchObj):
        replacements = {
            "\n": "\\n",
            "\r": "\\r",
            "\t": "\\t",
        }
        matchStr = matchObj.group(0)
        try:
            # If Matched portion of Regex matches a key
            # in the replacements dictionary, return replacement value
            # to be substituded; If key error, return the original matched string
            # (no substitution) and print a warning
            r = replacements[matchStr]
            return r
        except:
            print(
                f"Warning: Found match in replacement regex: {matchStr}; replacement value not availble; keeping string as-is"
            )
            return matchStr

    # Replace newline, tab, carriage return values with escaped version
    text = re.sub(r"\n|\r|\t", replacement, text)
    # Negative numbers need to be zero padded for valid json; replace -.* with -0.*
    text = re.sub(r"-\.", "-0.", text)
    # .replace(/elms\\null/g, 'elms\\"."') // hack to cover description field for LRB Location "Canaseraga Creek"

    return text


def cwms_office_symbols():
    """Returns an array of all office symbols in CWMS API"""
    offices_response = requests.get(
        "https://cwms-data.usace.army.mil/cwms-data/offices"
    ).json()

    return [f["name"] for f in offices_response["offices"]["offices"]]


def locations_for_office_symbols(office_symbols):
    """Returns a list of locations from CWMS RADAR API given an array of office symbols like ["lrh", "lrn"]"""
    locations = []
    for f in office_symbols:
        office_locations = requests.get(
            f"https://cwms-data.usace.army.mil/cwms-data/locations?&office={f}&names=@&format=json"
        )
        try:
            _j = office_locations.json()
            _j_locations = _j["locations"]["locations"]
            locations += _j_locations
            print(
                f"Successfully Extracted {len(_j_locations)} Locations for Office Symbol: {f}"
            )
        except:
            try:
                _j = radar_fix_json(office_locations.text)
                _j_locations = json.loads(_j)["locations"]["locations"]
                locations += _j_locations
            except:
                print(
                    f"ERROR: COULD NOT CONVERT API RESPONSE TO JSON FOR OFFICE SYMBOL: {f}"
                )

    return locations


# def office_id_by_symbol():
#     """Returns a dictionary with key of Office Symbol ("lrh", "lrn", "mvp", "nwo", etc.),
#     value of (UUID4). Based on "Shared API" used by multiple CWBI applications
#     """

#     url = "https://api.rsgis.dev/development/shared/offices"
#     offices = requests.get(url)
#     return {f["symbol"]: f["id"] for f in offices.json()}


if __name__ == "__main__":

    ###############
    # CONFIGURATION
    ###############
    # OFFICE_SYMBOLS = "ALL"
    OFFICE_SYMBOLS = ["lrn", "lrh", "mvp"]
    OUTFILE = "./seed_data.sql"
    ###############

    states = {
        "TN": 30,
        "MN": 4,
        "WV": 1,
        "GA": 22,
        "SC": 28,
        "NC": 9,
        "AL": 23,
        "NY": 39,
        "KY": 34,
        "VA": 53,
        "OH": 25,
    }

    print(
        f"Running Data Extraction for Specified Office Symbols: \n\t{OFFICE_SYMBOLS}\n"
    )

    print(f"Offices Available in CWMS RADAR API: \n\t{cwms_office_symbols()}\n")

    # Dict for reverse lookup of Office ID from Office Symbol
    # office_map = office_id_by_symbol()

    # Unknown offices; Keep track of office symbols used in CWMS that are not technically USACE Offices
    # Key: Office Symbol; Value: Number of Corresponding Locations
    unknown_offices_map = {}

    # Location Kind Map
    location_kind_map = {}

    # Known Slugs
    used_slugs = {}

    # List of CWMS Locations Retrieved from CWMS RADAR API
    if OFFICE_SYMBOLS == "ALL":
        locations = locations_for_office_symbols(cwms_office_symbols())
    else:
        locations = locations_for_office_symbols(OFFICE_SYMBOLS)

    _entries = []
    for idx, l in enumerate(locations):
        _entry = "("
        # Office ID
        cwms_office = l["identity"]["office"]
        # try:
        #     office_id = office_map[cwms_office]
        # except:
        #     if cwms_office not in unknown_offices_map:
        #         unknown_offices_map[cwms_office] = 1
        #     else:
        #         unknown_offices_map[cwms_office] += 1
        #     continue

        datasource_id_select = f"""(SELECT id FROM v_datasource  WHERE datatype = 'cwms-location' AND provider = '{cwms_office.lower()}')"""

        _entry += datasource_id_select + ","
        # Name

        _entry += quote(apostrophe_doubleup(l["identity"]["name"])) + ","
        # Public Name
        _public_name = l["label"]["public-name"]
        if _public_name is None:
            # No single quotes on null
            public_name = "null"
        else:
            public_name = '"' + _public_name + '"'
        # Slug
        slug = nextUniqueSlug(l["identity"]["name"], used_slugs)
        used_slugs[slug] = slug
        _entry += "slugify(" + quote(slug) + "),"
        # Geometry
        lon, lat = l["geolocation"]["longitude"], l["geolocation"]["latitude"]
        if lon is None or lat is None:
            lon, lat = 0, 0
        wkt = quote(f"POINT({lon} {lat})")
        _entry += f"ST_GeomFromText({wkt},4326),"
        # Kind ID
        location_kind = l["classification"]["location-kind"]
        if location_kind not in location_kind_map:
            location_kind_map[location_kind] = uuid.uuid4()

        _state = l["political"]["state"]
        try:
            state_id = str(states[_state.upper()]) + ","
        except:
            state_id = "NULL,"

        _entry += state_id

        # _entry += quote(location_kind_map[location_kind])
        attrbutes = f"""'{{"public_name": {public_name}, "kind": "{location_kind}"}}'"""
        _entry += attrbutes
        _entry += ")"

        _entries.append(_entry)

    with open(OUTFILE, "w") as f:
        # Insert new Location Kinds
        # f.write("INSERT INTO location_kind (id, name) VALUES\n")
        # for idx, k in enumerate(location_kind_map.keys()):
        #     _str = f"\t('{location_kind_map[k]}','{k}')"
        #     if idx != len(location_kind_map) - 1:
        #         _str += ",\n"
        #     else:
        #         _str += ";\n"
        #     f.write(_str)

        # Insert New Locations
        f.write("\n\n")
        f.write(
            f"INSERT INTO location (datasource_id, code, slug, geometry, state_id, attributes) VALUES\n"
        )
        for idx, entry in enumerate(_entries):
            if idx != len(_entries) - 1:
                f.write(f"\t{entry},\n")
            else:
                f.write(f"\t{entry};\n")

    f.close()

# Summary Report
for k in unknown_offices_map.keys():
    print(
        f"Office '{k}' is not recognized as a Corps office; {unknown_offices_map[k]} corresponding locations were not loaded;"
    )
