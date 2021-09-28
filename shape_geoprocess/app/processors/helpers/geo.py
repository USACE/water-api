import json
from geojson_rewind import rewind
import shapefile #pyshp

def shape_to_geojson(shape):
    
    schema = {                                                                                                     
        "geometry": "Polygon",                                                                                     
        "properties": {"id": "int"}                                                                                
    } 
    geojson_data = {}
    geojson_data['type'] = 'Feature'
    geojson_data['geometry'] = shape.__geo_interface__

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

    return geojson_data
