import json

from datetime import datetime, timezone

import boto3
import botocore
import botocore.exceptions

import config as CONFIG


INCOMING_SHAPEFILE_TO_SINGLE_SHAPE = {
    "processor": "watershed_shapefile_upload",
    "input": {
        "bucket": "castle-data-develop",
        "key": "water/test-watershed/LRH_Scioto.zip",
    },
    "functions": [
        {"function": "cleanup"},
        {"function": "dissolve"},
        {"function": "simplify"},
        {"function": "transform"},
    ],
    "output_target": "watersheds/5758d0dc-c8bf-4e37-a5e7-44ff3f4b8677/update_geometry",
}

CLIENT = boto3.resource(
    "sqs",
    endpoint_url="http://elasticmq:9324",
    region_name="elasticmq",
    aws_secret_access_key="x",
    aws_access_key_id="x",
    use_ssl=False,
)

# Incoming Requests
queue = CLIENT.get_queue_by_name(QueueName=CONFIG.QUEUE_NAME)

print(f"queue;       : {queue}")

msg = INCOMING_SHAPEFILE_TO_SINGLE_SHAPE

response = queue.send_message(MessageBody=json.dumps(msg, separators=(",", ":")))
