"""
Notes:
The following is used by GDAL (passed to Fiona) to override 
S3 defaults and allow Minio local testing:

AWS_S3_ENDPOINT
AWS_VIRTUAL_HOSTING=FALSE
AWS_HTTPS=NO
"""

import os

APPLICATION_KEY = os.getenv("APPLICATION_KEY", default="appkey")
ENDPOINT_URL_SQS = os.getenv("ENDPOINT_URL_SQS", default=None)
QUEUE_NAME = os.getenv("QUEUE_NAME", default="water-geoprocess")

#
# AWS Credentials
#
AWS_SECRET_ACCESS_KEY = os.getenv("AWS_SECRET_ACCESS_KEY", default="x")
AWS_ACCESS_KEY_ID = os.getenv("AWS_ACCESS_KEY_ID", default="x")
AWS_REGION = os.getenv("AWS_REGION", default="us-east-1")

# If _SQS versions of AWS_SECRET_ACCESS_KEY, AWS_ACCESS_KEY_ID, AWS_REGION not explicitly set,
# set to usual environment variables for AWS credentials. Variables with _SQS suffix are intended
# to be used when override is required for elasticmq (variables set to 'x' or 'elasticmq') for local testing.
# See Documentation Here: https://github.com/softwaremill/elasticmq#using-the-amazon-boto-python-to-access-an-elasticmq-server
AWS_SECRET_ACCESS_KEY_SQS = os.getenv(
    "AWS_SECRET_ACCESS_KEY_SQS", default=AWS_SECRET_ACCESS_KEY
)
AWS_ACCESS_KEY_ID_SQS = os.getenv("AWS_ACCESS_KEY_ID_SQS", default=AWS_ACCESS_KEY_ID)
AWS_REGION_SQS = os.getenv("AWS_REGION_SQS", default=AWS_REGION)

#
# Configuration Parameters
#
USE_SSL = os.getenv("USE_SSL", default=False)

WAIT_TIME_SECONDS = os.getenv("WAIT_TIME_SECONDS", default=20)

WRITE_TO_BUCKET = os.getenv("WRITE_TO_BUCKET", default="castle-data-develop")

WATER_API_URL = os.getenv("WATER_API_URL", default="http://api")
WATER_API_HOST_HEADER = os.getenv("WATER_API_HOST_HEADER", default=None)
