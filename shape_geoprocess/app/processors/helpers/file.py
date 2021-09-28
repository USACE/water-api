import json
import fiona
from shapely.geometry import shape

'''
Returns a list of shapely geometry objects
'''
def read_s3_zip(s3_path):
    with fiona.open(f'zip+s3://{s3_path}', 'r') as collection:

        shapes = []
        for x in list(collection):
            shapes.append(shape(x["geometry"]))        
        crs = collection.crs
    return (shapes, crs)

# src = "cwbi-data-develop/water/test-watershed/LRH_Scioto.zip"
# shapes = read_s3_zip(src)
# print(shapes)

def write_geojson(geojson_dict, fileout):
    # Write to geojson file
    with open(fileout, 'w') as outfile:
        outfile.write(json.dumps(geojson_dict, indent = 4))