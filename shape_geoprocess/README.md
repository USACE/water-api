# SHAPE_GEOPROCESS

## About:
**shape_geoprocess** will accept an incoming SQS message and perform ths following:
1) Read the shapefile zip on S3
2) Perform the processing defined in the SQS `processes` list
3) Convert shape to GeoJSON
4) Post to API defined in `output_target`

## Sample SQS Body Payload
```json
{
    "processor": "watershed_shapefile_upload",
    "input": {
        "bucket": "cwbi-data-develop",
        "key": "water/test-watershed/LRH_Scioto.zip",
    },
    "functions": [
        {"function": "cleanup"},
        {"function": "dissolve"},
        {"function": "simplify"},
        {"function": "transform"},
    ],
    "output_target": "http://api/watersheds/geometry/5758d0dc-c8bf-4e37-a5e7-44ff3f4b8677",
}
```

## To Start:

1) Ensure ElasticMQ is running (should be in main docker-compose)
2) Ensure minio is running `./minio.sh start`
3) `docker-compose up` in shape_geoprocess
4) Jump into shape_geoprocess container and run `python test-message.py`