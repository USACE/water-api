import logging
import json
import importlib

import requests
import boto3
import botocore
import botocore.exceptions

import config as CONFIG

from processors.helpers.file import read_s3_zip, write_geojson
from processors.helpers.geo import shape_to_geojson

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
print(f'queue; shape-geoprocess       : {queue_packager}')

# Logger
logger = logging.getLogger()
logger.setLevel(logging.INFO)
logger.addHandler(logging.StreamHandler())

def get_infile_processor(name):
    """Import library for processing a given product_name"""
    
    processor = importlib.import_module(f'processors.{name}')
    
    return processor

def handle_message(msg):
    """Converts JSON-Formatted message string to dictionary and calls package()"""

    print('\n\nmessage received\n\n')
    print(msg.body)

    body = json.loads(msg.body)
    bucket, key = body['input']['bucket'], body['input']['key']
    # watershed_id = body['input']['watershed_id']
    processes = body['processes']
    output_target = body['output_target']
    
    # Read in S3 Zip, load into shapes object
    try:
        shapes, crs = read_s3_zip(f'{bucket}/{key}')
    except:
        logger.error(f'file not found: {bucket}/{key}')
        return
    
    # Process based on list of processes in sqs message
    for p in processes:
        processor = get_infile_processor(p['process'])
        logger.info(f'Using processor: {processor}')
        if p['process'] == 'transform':
            shapes = processor.process(shapes, crs)
        else:
            shapes = processor.process(shapes)

    # Convert the shape object to geojson
    logger.info('Converting shape to GeoJSON')
    geojson = shape_to_geojson(shapes)

    write_geojson(geojson, '/usr/src/app/testoutput.json')

    logger.info(f'Sending GeoJson to API target: {output_target}')
    r = requests.put(f'{output_target}?key={CONFIG.APPLICATION_KEY}', 
    json=geojson)
    if r.status_code != 200:
        logger.info(r.status_code)
        logger.info(r.text)    
    
    

while 1:
    messages = queue_packager.receive_messages(WaitTimeSeconds=CONFIG.WAIT_TIME_SECONDS)
    print(f'shapeprocessor message count: {len(messages)}')
    
    for message in messages:
        handle_message(message)
        message.delete()