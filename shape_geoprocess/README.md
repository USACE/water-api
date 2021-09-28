# SHAPE_GEOPROCESS

## To Start:

1) Ensure ElasticMQ is running (should be in main docker-compose)
2) Ensure minio is running `./minio.sh start`
3) `docker-compose up` in shape_geoprocess
4) Jump into shape_geoprocess container and run `python test-message.py`