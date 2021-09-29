import json
import logging
import importlib

import requests

from .helpers.file import read_s3_zip, write_geojson

import config as CONFIG

from .shape_functions import shape_to_geojson


# Logger
logger = logging.getLogger()


def get_function(name):
    """Import library for processing a given product_name"""

    return importlib.import_module(f"processors.shape_functions.{name}")    

def process(body):

    # print(body['input']['bucket'])
    path = f"{body['input']['bucket']}/{body['input']['key']}"
    output_target = body['output_target']
    
    # Read in S3 Zip, load into shapes object
    try:
        shapes, crs = read_s3_zip(path)
    except:
        logger.error(f'file not found: {path}')
        return

    for f in body['functions']:
        fn_name = f['function']
        
        function = get_function(fn_name)

        if fn_name == 'transform':
            shapes = function.process(shapes, crs)
        else:
            shapes = function.process(shapes)

    # Convert the shape object to geojson
    logger.info('Converting shape to GeoJSON')
    geojson = shape_to_geojson(shapes)

    # write_geojson(geojson, '/usr/src/app/testoutput.json')

    logger.info(f'Sending GeoJson to API target: {output_target}')
    r = requests.put(f'{output_target}?key={CONFIG.APPLICATION_KEY}', 
    json=geojson)
    if r.status_code != 200:
        logger.info(r.status_code)
        logger.info(r.text)    