import os
import json
import logging

import boto3
import botocore
import botocore.exceptions
from fiona.session import AWSSession
import fiona                                                                                                       
from shapely.ops import unary_union, transform                                                                           
from shapely.geometry import shape, mapping  
from shapely.validation import make_valid
from fiona.crs import to_string
import shapefile
from geojson_rewind import rewind
from pyproj import Proj, Transformer, CRS
import requests

import config as CONFIG

if CONFIG.AWS_ACCESS_KEY_ID == 'x':
    # Running in AWS
    # Using IAM Role for Credentials
    CLIENT = boto3.resource('sqs')
else:
    # Local Testing
    # ElasticMQ with Credentials via AWS_ environment variables
    CLIENT = boto3.resource(
        'sqs',
        endpoint_url=CONFIG.ENDPOINT_URL_SQS,
        region_name=CONFIG.AWS_REGION_SQS,
        aws_secret_access_key=CONFIG.AWS_SECRET_ACCESS_KEY_SQS,
        aws_access_key_id=CONFIG.AWS_ACCESS_KEY_ID_SQS,
        use_ssl=CONFIG.USE_SSL
    )

# Incoming Requests for Processor
queue_packager = CLIENT.get_queue_by_name(QueueName=CONFIG.QUEUE_NAME)
print(f'queue; shape-processor       : {queue_packager}')

# Logger
logger = logging.getLogger()
logger.setLevel(logging.INFO)
logger.addHandler(logging.StreamHandler())

def get_infile(bucket, key, filepath):
    
    s3 = boto3.client('s3')
    try:
        s3.Bucket(bucket).download_file(key, filepath)
        return os.path.abspath(filepath)

    except botocore.exceptions.ClientError as e:
        if e.response['Error']['Code'] == "404":
            logger.fatal(f'OBJECT DOES NOT EXIST: {key}')
            return None
        else:
            raise

def handle_message(msg):
    """Converts JSON-Formatted message string to dictionary and calls package()"""

    print('\n\nmessage received\n\n')
    print(msg.body)

    print(type(msg.body))

    body = json.loads(msg.body)
    bucket = body['process_config']['bucket']
    key = body['process_config']['key']
    watershed_id = body['process_config']['watershed_id']

    src = f"zip+s3://{bucket}/{key}"

    

    with fiona.open(src, 'r') as ds_in:                                                                                                                                                                                                                   
        crs = ds_in.crs 
        dst_crs = fiona.crs.from_epsg(4326)
        drv = ds_in.driver 
    

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

    r = requests.put(f'http://api/watersheds/geometry/{watershed_id}?key=appkey', 
        json=geojson_data)
    print(r.status_code)
    print(r.text)


    # watershed_id = json.loads(msg.body)["watershed_id"]
    # r = requests.get(
    #     f'{CONFIG.CUMULUS_API_URL}/downloads/{download_id}/packager_request',
    # )
    # if r.status_code == 200:
    #     package(r.json(), packager_update_fn)
    # else:
    #     brk = '*' * 24
    #     print(f'{brk}\nPackager Fail On Message: {msg}{brk}\n')
    #     print(f'{brk}\nRequest: {r.request.url}{brk}\n')
    #     print(f'{brk}\nHeaders: {r.request.headers}{brk}\n')
    #     print(f'{brk}\nStatus Code: {r.status_code}{brk}\n')
    #     print(f'{brk}\nReason: {r.reason}{brk}\n')
    #     print(f'{brk}\nContent: {r.content}{brk}\n')


while 1:
    messages = queue_packager.receive_messages(WaitTimeSeconds=CONFIG.WAIT_TIME_SECONDS)
    print(f'packager message count: {len(messages)}')
    
    for message in messages:
        handle_message(message)
        message.delete()