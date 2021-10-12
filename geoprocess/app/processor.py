import logging
import json
import importlib

import boto3
import botocore
import botocore.exceptions

import config as CONFIG

if CONFIG.AWS_ACCESS_KEY_ID == "x":
    # Running in AWS
    # Using IAM Role for Credentials
    CLIENT = boto3.resource("sqs")
else:
    # Local Testing
    # ElasticMQ with Credentials via AWS_ environment variables
    CLIENT = boto3.resource(
        "sqs",
        endpoint_url=CONFIG.ENDPOINT_URL_SQS,
        region_name=CONFIG.AWS_REGION_SQS,
        aws_secret_access_key=CONFIG.AWS_SECRET_ACCESS_KEY_SQS,
        aws_access_key_id=CONFIG.AWS_ACCESS_KEY_ID_SQS,
        use_ssl=CONFIG.USE_SSL,
    )

# Incoming Requests for Processor
queue_packager = CLIENT.get_queue_by_name(QueueName=CONFIG.QUEUE_NAME)
print(f"queue; shape-geoprocess       : {queue_packager}")

# Logger
logger = logging.getLogger()
logger.setLevel(logging.INFO)
logger.addHandler(logging.StreamHandler())


def get_processor(name):
    """Import library for processing a given product_name"""

    processor = importlib.import_module(f"processors.{name}")
    return processor


def handle_message(msg):
    """Converts JSON-Formatted message string to dictionary and calls package()"""

    print("\n\nmessage received\n\n")
    print(msg.body)
    body = json.loads(msg.body)

    # Process based on list of processes in sqs message
    processor = get_processor(body["processor"])

    logger.info(f"Using processor: {processor}")
    processor.process(body)


while 1:
    messages = queue_packager.receive_messages(WaitTimeSeconds=CONFIG.WAIT_TIME_SECONDS)
    print(f"shapeprocessor message count: {len(messages)}")

    for message in messages:
        handle_message(message)
        message.delete()
