import json
import fiona                                                                                                       
from shapely.ops import unary_union, transform                                                                           
from shapely.geometry import shape, mapping  
from shapely.validation import make_valid
from fiona.crs import to_string
import shapefile
from geojson_rewind import rewind
from pyproj import Proj, Transformer, CRS

src = 'zip:///test_shapes/LRH_Scioto.zip'
dst = '/test_shapes/fixed_LRH_Scioto.shp'


with fiona.open(src, 'r') as ds_in:                                                                                                                                                                                                                   
    crs = ds_in.crs 
    dst_crs = fiona.crs.from_epsg(4326)
    drv = ds_in.driver  
   
        
    # print(ds_in.schema)
    # print(f'CRS is: {crs}')
    # print(to_string(ds_in.crs))

    shapes = []
    invalid_shapes = []
    for x in list(ds_in):
        if not shape(x["geometry"]).is_valid:
            # Shape is not valid, apply fix
            invalid_shapes.append(shape)
            valid = make_valid(shape(x["geometry"])).buffer(0.5)
            shapes.append(valid)            
        else:
            shapes.append(shape(x["geometry"]).buffer(0.5))

    print(f'{len(invalid_shapes)} of {len(shapes)} shapes were invalid and corrected.')                                                                                 
                                              
    # Dissolve multiple shapes into single shape
    dissolved_shape = unary_union(shapes)           

schema = {                                                                                                     
    "geometry": "Polygon",                                                                                     
    "properties": {"id": "int"}                                                                                
}  

print(len(dissolved_shape.exterior.coords))
'''
# https://shapely.readthedocs.io/en/stable/manual.html#object.simplify
All points in the simplified object will be within the tolerance distance of the 
original geometry. By default a slower algorithm is used that preserves topology. 
If preserve topology is set to False the much quicker Douglas-Peucker algorithm is used.
'''
simple_shape = dissolved_shape.simplify(0.5, preserve_topology=False)

# Credit: https://gis.stackexchange.com/questions/127427/transforming-shapely-polygon-and-multipolygon-objects
project = Transformer.from_proj(
    Proj(crs), # source coordinate system
    Proj(CRS('EPSG:4326')), # destination coordinate system. Note: not using dst_crs as pyproj will complain 
                            # with future warning about soon to be unsupported authority/code string
    always_xy=True #True keeps the lon first in the coordinate
    ) 

# Transform all points in the polygon from src_crs to dst_crs
simple_shape = transform(project.transform, simple_shape)  # apply projection


geojson_data = {}
geojson_data['type'] = 'Feature'
geojson_data['geometry'] = simple_shape.__geo_interface__

# print(geojson_data['geometry']['coordinates'])

'''
https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.6
A linear ring MUST follow the right-hand rule with respect to the
area it bounds, i.e., exterior rings are counterclockwise, and
holes are clockwise.
'''
# Rewind will perform the right-hand winding
geojson_data_string = rewind(json.dumps(geojson_data))
geojson_data = json.loads(geojson_data_string)
geojson_data['properties'] = None

# Write to geojson file
with open('/test_shapes/LRH_Scioto.json', 'w') as outfile:
    outfile.write(json.dumps(geojson_data, indent = 4))

# Write to "fixed" shapefile
with fiona.open(dst, 'w', driver=drv, schema=schema, crs=dst_crs) as ds_dst:
    ds_dst.write({"geometry": mapping(simple_shape), "properties": {"id": 1}})